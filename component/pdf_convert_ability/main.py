import argparse
import json
from convert import pdf_convert

def main():
    parser = argparse.ArgumentParser(description='PDF to Markdown converter')
    parser.add_argument("--pdf_path", required=True, help="Path to the input PDF file")
    parser.add_argument("--output_dir", required=True, help="Directory to save the output files")

    args = parser.parse_args()
    try:
        md_path, md_dir = pdf_convert(args.pdf_path, args.output_dir, True)
        # 返回 JSON 格式，方便其他语言解析
        print(json.dumps({"success": True, "md_path": md_path, "md_dir": md_dir}))
    except Exception as e:
        print(json.dumps({"success": False, "error": str(e)}))
        exit(1)  # 非零状态表示执行失败

if __name__ == "__main__":
    main()
