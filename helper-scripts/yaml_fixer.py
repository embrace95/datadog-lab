import os
import sys
import difflib
import argparse
from ruamel.yaml import YAML
from ruamel.yaml.parser import ParserError
from ruamel.yaml.scanner import ScannerError

yaml = YAML()
yaml.indent(mapping=2, sequence=4, offset=2)
yaml.preserve_quotes = True

def fix_yaml(file_path, auto_overwrite=False):
    try:
        with open(file_path, 'r', encoding='utf-8') as f:
            original = f.read()
        data = yaml.load(original)

        from io import StringIO
        fixed_stream = StringIO()
        yaml.dump(data, fixed_stream)
        fixed = fixed_stream.getvalue()

        if original != fixed:
            print(f"\n📄 Changes detected in {file_path}")
            for line in difflib.unified_diff(
                original.splitlines(keepends=True),
                fixed.splitlines(keepends=True),
                fromfile='original',
                tofile='fixed'
            ):
                sys.stdout.write(line)

            if auto_overwrite:
                with open(file_path, 'w', encoding='utf-8') as f:
                    f.write(fixed)
                print(f"✅ Overwritten: {file_path}")
            else:
                print("⚠️ Not overwritten. Use --write to apply changes.")

        else:
            print(f"✅ No changes needed in {file_path}")

    except (ParserError, ScannerError) as e:
        print(f"\n❌ YAML error in {file_path}:\n{e}")

def find_yaml_files(path):
    for root, _, files in os.walk(path):
        for file in files:
            if file.endswith((".yaml", ".yml")):
                yield os.path.join(root, file)

if __name__ == "__main__":
    parser = argparse.ArgumentParser(description="Fix YAML formatting with diff and syntax logging.")
    parser.add_argument("target", help="Path to YAML file or folder")
    parser.add_argument("--write", action="store_true", help="Automatically overwrite files with fixes")
    args = parser.parse_args()

    target = args.target
    if os.path.isfile(target):
        fix_yaml(target, auto_overwrite=args.write)
    elif os.path.isdir(target):
        for yaml_file in find_yaml_files(target):
            fix_yaml(yaml_file, auto_overwrite=args.write)
    else:
        print("❌ Invalid path: must be file or directory")

