import os


def create_md(document_path:str, output_dir:str) -> tuple[str, str]:
    file_name = os.path.basename(document_path)
    file_name_only, file_ext = os.path.splitext(file_name)
    md_name = file_name_only+".md"
    md_path = os.path.join(output_dir, md_name)
    md_dir_path = output_dir
    # 检查md_dir_path是否存在，不存在则 创建
    if not os.path.exists(md_dir_path):
        os.makedirs(md_dir_path)
    if not os.path.exists(os.path.join(md_dir_path, "images")):
        os.makedirs(os.path.join(md_dir_path, "images"))
    return md_path, md_dir_path
