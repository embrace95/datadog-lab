// Copyright 2018 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"time"


	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/GoogleCloudPlatform/microservices-demo/src/checkoutservice/genproto"
	money "github.com/GoogleCloudPlatform/microservices-demo/src/checkoutservice/money"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"

    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
    "go.opentelemetry.io/otel/exporters/otlp/otlpmetric"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
    "go.opentelemetry.io/otel/propagation"
    "go.opentelemetry.io/otel/attribute"

    "go.opentelemetry.io/otel/metric/global"
    "go.opentelemetry.io/otel/sdk/resource"
    "go.opentelemetry.io/otel/trace"
    "go.opentelemetry.io/otel/metric/instrument"
    "go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
    sdktrace "go.opentelemetry.io/otel/sdk/trace"
    semconv  "go.opentelemetry.io/otel/semconv/v1.7.0"
    controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"
    processor "go.opentelemetry.io/otel/sdk/metric/processor/basic"
    "go.opentelemetry.io/otel/sdk/metric/selector/simple"
)

const (
	listenPort  = "5050"
	usdCurrency = "USD"
)
var tracer = otel.Tracer("checkoutService-tracer")
var meter = global.MeterProvider().Meter("checkoutService-meter")
var log *logrus.Logger

func init() {
	log = logrus.New()
	log.Level = logrus.DebugLevel
	log.Formatter = &logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "timestamp",
			logrus.FieldKeyLevel: "severity",
			logrus.FieldKeyMsg:   "message",
		},
		TimestampFormat: time.RFC3339Nano,
	}
	log.Out = os.Stdout
}

type checkoutService struct {
	productCatalogSvcAddr string
	cartSvcAddr           string
	currencySvcAddr       string
	shippingSvcAddr       string
	emailSvcAddr          string
	paymentSvcAddr        string
}
func  initProvider()  {
	ctx := context.Background()

    	res, err := resource.New(ctx,
        		resource.WithFromEnv(),
        		resource.WithProcess(),
        		resource.WithTelemetrySDK(),
        		resource.WithHost(),
        		resource.WithAttributes(
        			// the service name used to display traces in backends
        			semconv.ServiceNameKey.String("checkout-Service"),
        		),
        	)
        	handleErr(err, "failed to create resource")

    	// If the OpenTelemetry Collector is running on a local cluster (minikube or
    	// microk8s), it should be accessible through the NodePort service at the
    	// `localhost:30080` endpoint. Otherwise, replace `localhost` with the
    	// endpoint of your cluster. If you run the app inside k8s, then you can
    	// probably connect directly to the service through dns
        svcAddr := os.Getenv("OTLP_SERVICE_ADDR")
        if svcAddr == "" {
            log.Info("OpenTelemetry initialization disabled.")
            return
        }
        svcPort := os.Getenv("OTLP_SERVICE_PORT")
                if svcPort == "" {
                    log.Info("OpenTelemetry initialization disabled.")
                    return
                }



    	conn, err := grpc.DialContext(ctx, svcAddr+":"+svcPort, grpc.WithInsecure(), grpc.WithBlock())
        handleErr(err, "failed to create gRPC connection to collector")

         //initialization of the trace provider
        otlpTraceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
        handleErr(err, "failed to create trace exporter")

        batchSpanProcessor := sdktrace.NewBatchSpanProcessor(otlpTraceExporter)

        tracerProvider := sdktrace.NewTracerProvider(
                sdktrace.WithSampler(sdktrace.AlwaysSample()),
                sdktrace.WithSpanProcessor(batchSpanProcessor),
                sdktrace.WithResource(res),
        )

       /// init of the metric exporter
        metricClient := otlpmetricgrpc.NewClient(
                   otlpmetricgrpc.WithInsecure(),
                   otlpmetricgrpc.WithGRPCConn(conn))

       metricExp, err := otlpmetric.New(ctx, metricClient)
       handleErr(err, "Failed to create the collector metric exporter")
        //Create the metric provider
       pusher := controller.New(
           processor.NewFactory(
               simple.NewWithHistogramDistribution(),
               metricExp,
           ),
           controller.WithExporter(metricExp),
           controller.WithCollectPeriod(2*time.Second),
       )
        global.SetMeterProvider(pusher)
        err = pusher.Start(ctx)
        handleErr(err, "Failed to start metric pusher")
        otel.SetTracerProvider(tracerProvider)
        otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
        propagation.TraceContext{}, propagation.Baggage{}))
}
func main() {

    initProvider()

	port := listenPort
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}




	svc := new(checkoutService)
	mustMapEnv(&svc.shippingSvcAddr, "SHIPPING_SERVICE_ADDR")
	mustMapEnv(&svc.productCatalogSvcAddr, "PRODUCT_CATALOG_SERVICE_ADDR")
	mustMapEnv(&svc.cartSvcAddr, "CART_SERVICE_ADDR")
	mustMapEnv(&svc.currencySvcAddr, "CURRENCY_SERVICE_ADDR")
	mustMapEnv(&svc.emailSvcAddr, "EMAIL_SERVICE_ADDR")
	mustMapEnv(&svc.paymentSvcAddr, "PAYMENT_SERVICE_ADDR")

	log.Infof("service config: %+v", svc)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatal(err)
	}

	var srv *grpc.Server
	if os.Getenv("DISABLE_STATS") == "" {
		log.Info("Stats enabled.")
		srv = grpc.NewServer(
		        grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()),
                grpc.StreamInterceptor(otelgrpc.StreamServerInterceptor()),
        		)
	} else {
		log.Info("Stats disabled.")
		srv = grpc.NewServer(
                grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()),
                grpc.StreamInterceptor(otelgrpc.StreamServerInterceptor()),
                )
	}
	pb.RegisterCheckoutServiceServer(srv, svc)
	healthpb.RegisterHealthServer(srv, svc)
	log.Infof("starting to listen on tcp: %q", lis.Addr().String())
	err = srv.Serve(lis)
	log.Fatal(err)
}



func mustMapEnv(target *string, envKey string) {
	v := os.Getenv(envKey)
	if v == "" {
		panic(fmt.Sprintf("environment variable %q not set", envKey))
	}
	*target = v
}

func (cs *checkoutService) Check(ctx context.Context, req *healthpb.HealthCheckRequest) (*healthpb.HealthCheckResponse, error) {
	return &healthpb.HealthCheckResponse{Status: healthpb.HealthCheckResponse_SERVING}, nil
}

func (cs *checkoutService) Watch(req *healthpb.HealthCheckRequest, ws healthpb.Health_WatchServer) error {
	return status.Errorf(codes.Unimplemented, "health check via Watch not implemented")
}

func (cs *checkoutService) PlaceOrder(ctx context.Context, req *pb.PlaceOrderRequest) (*pb.PlaceOrderResponse, error) {
	log.Infof("[PlaceOrder] user_id=%q user_currency=%q", req.UserId, req.UserCurrency)

    newCtx, span := tracer.Start(ctx, "PlaceOrder")




    ordercount,err := meter.SyncInt64().
		Counter(
			"hipstershop.checkoutService.order_counts",
			instrument.WithDescription("The counts of orders"),

		)
	if err != nil {
    		log.Errorf("Failed to register instrument")
    	}

    orderitemsCount,err :=  meter.SyncInt64().
		Counter(
			"hipstershop.checkoutService.order_items_counts",
			instrument.WithDescription("The counts of items ordered"),

		)
    if err != nil {
    		log.Errorf("Failed to register instrument")
    }
	orderID, err := uuid.NewUUID()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate order uuid")
	}

	prep, err := cs.prepareOrderItemsAndShippingQuoteFromCart(ctx, req.UserId, req.UserCurrency, req.Address)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	total := pb.Money{CurrencyCode: req.UserCurrency,
		Units: 0,
		Nanos: 0}
	total = money.Must(money.Sum(total, *prep.shippingCostLocalized))
	var items int32
	for _, it := range prep.orderItems {
		multPrice := money.MultiplySlow(*it.Cost, uint32(it.GetItem().GetQuantity()))
		items+=it.GetItem().GetQuantity()
		total = money.Must(money.Sum(total, multPrice))
	}


	txID, err := cs.chargeCard(ctx, &total, req.CreditCard)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to charge card: %+v", err)
	}
	log.Infof("payment went through (transaction_id: %s)", txID)

	shippingTrackingID, err := cs.shipOrder(ctx, req.Address, prep.cartItems)
	if err != nil {
		return nil, status.Errorf(codes.Unavailable, "shipping error: %+v", err)
	}

	_ = cs.emptyUserCart(ctx, req.UserId)

	orderResult := &pb.OrderResult{
		OrderId:            orderID.String(),
		ShippingTrackingId: shippingTrackingID,
		ShippingCost:       prep.shippingCostLocalized,
		ShippingAddress:    req.Address,
		Items:              prep.orderItems,
	}

	if err := cs.sendOrderConfirmation(ctx, req.Email, orderResult); err != nil {
		log.Warnf("failed to send order confirmation to %q: %+v", req.Email, err)
	} else {
		log.Infof("order confirmation email sent to %q", req.Email)
	}
	resp := &pb.PlaceOrderResponse{Order: orderResult}

    ordercount.Add(newCtx,1,attribute.String("method", "PlaceOrder"),
                                		attribute.String("client", "Checkoutservice"),
                                		attribute.String("user_id", req.UserId),
                                    	attribute.String("user_currency", req.UserCurrency))
    orderitemsCount.Add(newCtx,int64(items),attribute.String("method", "PlaceOrder"),
                                                		attribute.String("client", "Checkoutservice"),
                                                		attribute.String("user_id", req.UserId),
                                                    	attribute.String("user_currency", req.UserCurrency))


    span.AddEvent("Place ordered", trace.WithAttributes(
      attribute.Int("number of items",int(items)),
    ))

	span.End()
	return resp, nil
}

type orderPrep struct {
	orderItems            []*pb.OrderItem
	cartItems             []*pb.CartItem
	shippingCostLocalized *pb.Money
}

func (cs *checkoutService) prepareOrderItemsAndShippingQuoteFromCart(ctx context.Context, userID, userCurrency string, address *pb.Address) (orderPrep, error) {
	var out orderPrep
	cartItems, err := cs.getUserCart(ctx, userID)
	if err != nil {
		return out, fmt.Errorf("cart failure: %+v", err)
	}
	orderItems, err := cs.prepOrderItems(ctx, cartItems, userCurrency)
	if err != nil {
		return out, fmt.Errorf("failed to prepare order: %+v", err)
	}
	shippingUSD, err := cs.quoteShipping(ctx, address, cartItems)
	if err != nil {
		return out, fmt.Errorf("shipping quote failure: %+v", err)
	}
	shippingPrice, err := cs.convertCurrency(ctx, shippingUSD, userCurrency)
	if err != nil {
		return out, fmt.Errorf("failed to convert shipping cost to currency: %+v", err)
	}

	out.shippingCostLocalized = shippingPrice
	out.cartItems = cartItems
	out.orderItems = orderItems
	return out, nil
}

func (cs *checkoutService) quoteShipping(ctx context.Context, address *pb.Address, items []*pb.CartItem) (*pb.Money, error) {
	conn, err := grpc.Dial( cs.shippingSvcAddr,
		grpc.WithInsecure(),
		  grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
                              grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()),
                        )
	if err != nil {
		return nil, fmt.Errorf("could not connect shipping service: %+v", err)
	}
	defer conn.Close()

	shippingQuote, err := pb.NewShippingServiceClient(conn).
		GetQuote(ctx, &pb.GetQuoteRequest{
			Address: address,
			Items:   items})
	if err != nil {
		return nil, fmt.Errorf("failed to get shipping quote: %+v", err)
	}
	return shippingQuote.GetCostUsd(), nil
}

func (cs *checkoutService) getUserCart(ctx context.Context, userID string) ([]*pb.CartItem, error) {
	conn, err := grpc.Dial( cs.cartSvcAddr,
	grpc.WithInsecure(),   grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
                                               grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()),)
	if err != nil {
		return nil, fmt.Errorf("could not connect cart service: %+v", err)
	}
	defer conn.Close()

	cart, err := pb.NewCartServiceClient(conn).GetCart(ctx, &pb.GetCartRequest{UserId: userID})
	if err != nil {
		return nil, fmt.Errorf("failed to get user cart during checkout: %+v", err)
	}
	return cart.GetItems(), nil
}

func (cs *checkoutService) emptyUserCart(ctx context.Context, userID string) error {
	conn, err := grpc.Dial(cs.cartSvcAddr, grpc.WithInsecure(),
	  grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
                          grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()),)
	if err != nil {
		return fmt.Errorf("could not connect cart service: %+v", err)
	}
	defer conn.Close()

	if _, err = pb.NewCartServiceClient(conn).EmptyCart(ctx, &pb.EmptyCartRequest{UserId: userID}); err != nil {
		return fmt.Errorf("failed to empty user cart during checkout: %+v", err)
	}
	return nil
}

func (cs *checkoutService) prepOrderItems(ctx context.Context, items []*pb.CartItem, userCurrency string) ([]*pb.OrderItem, error) {
	out := make([]*pb.OrderItem, len(items))

	conn, err := grpc.Dial( cs.productCatalogSvcAddr,
	 grpc.WithInsecure(),  grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
                                               grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()),)
	if err != nil {
		return nil, fmt.Errorf("could not connect product catalog service: %+v", err)
	}
	defer conn.Close()
	cl := pb.NewProductCatalogServiceClient(conn)

	for i, item := range items {
		product, err := cl.GetProduct(ctx, &pb.GetProductRequest{Id: item.GetProductId()})
		if err != nil {
			return nil, fmt.Errorf("failed to get product #%q", item.GetProductId())
		}
		price, err := cs.convertCurrency(ctx, product.GetPriceUsd(), userCurrency)
		if err != nil {
			return nil, fmt.Errorf("failed to convert price of %q to %s", item.GetProductId(), userCurrency)
		}
		out[i] = &pb.OrderItem{
			Item: item,
			Cost: price}
	}
	return out, nil
}

func (cs *checkoutService) convertCurrency(ctx context.Context, from *pb.Money, toCurrency string) (*pb.Money, error) {
	conn, err := grpc.Dial(cs.currencySvcAddr,
	 grpc.WithInsecure(),   grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
                                                grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()),)
	if err != nil {
		return nil, fmt.Errorf("could not connect currency service: %+v", err)
	}
	defer conn.Close()
	result, err := pb.NewCurrencyServiceClient(conn).Convert(context.TODO(), &pb.CurrencyConversionRequest{
		From:   from,
		ToCode: toCurrency})
	if err != nil {
		return nil, fmt.Errorf("failed to convert currency: %+v", err)
	}
	return result, err
}

func (cs *checkoutService) chargeCard(ctx context.Context, amount *pb.Money, paymentInfo *pb.CreditCardInfo) (string, error) {
	conn, err := grpc.Dial( cs.paymentSvcAddr,
	 grpc.WithInsecure(),   grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
                                                grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()),)
	if err != nil {
		return "", fmt.Errorf("failed to connect payment service: %+v", err)
	}
	defer conn.Close()

	paymentResp, err := pb.NewPaymentServiceClient(conn).Charge(ctx, &pb.ChargeRequest{
		Amount:     amount,
		CreditCard: paymentInfo})
	if err != nil {
		return "", fmt.Errorf("could not charge the card: %+v", err)
	}
	return paymentResp.GetTransactionId(), nil
}

func (cs *checkoutService) sendOrderConfirmation(ctx context.Context, email string, order *pb.OrderResult) error {
	conn, err := grpc.Dial( cs.emailSvcAddr,
	 grpc.WithInsecure(),
	   grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
                           grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()),)
	if err != nil {
		return fmt.Errorf("failed to connect email service: %+v", err)
	}
	defer conn.Close()
	_, err = pb.NewEmailServiceClient(conn).SendOrderConfirmation(ctx, &pb.SendOrderConfirmationRequest{
		Email: email,
		Order: order})
	return err
}

func (cs *checkoutService) shipOrder(ctx context.Context, address *pb.Address, items []*pb.CartItem) (string, error) {
	conn, err := grpc.Dial(cs.shippingSvcAddr,
	grpc.WithInsecure(),
	  grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
                    grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()),)
	if err != nil {
		return "", fmt.Errorf("failed to connect email service: %+v", err)
	}
	defer conn.Close()
	resp, err := pb.NewShippingServiceClient(conn).ShipOrder(ctx, &pb.ShipOrderRequest{
		Address: address,
		Items:   items})
	if err != nil {
		return "", fmt.Errorf("shipment failed: %+v", err)
	}
	return resp.GetTrackingId(), nil
}
func handleErr(err error, message string) {
	if err != nil {
		log.Fatalf("%s: %v",  message, err)
    }
}
