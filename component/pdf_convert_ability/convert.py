import os
from pdf_parser.pdf_reader import read_as_dict
from pdf_parser.block_types import BlockType
from pdf_parser.create_md import create_md
from pdf_parser.utils import save_base64_image, word2heading


def pdf_convert(pdf_path: str, output_dir: str, with_page_line:bool) -> tuple[str, str]:
    page_content_list = []
    md_path, md_dir = create_md(pdf_path, output_dir)
    page_list = read_as_dict(pdf_path)
    for i, page in enumerate(page_list):
        contents = []
        image_code = 1
        for block in page['chunks']:
            # block['font_size'] 映射表
            if block['type'] == BlockType.TEXT.value:
                text = word2heading(block['text'].strip(), block['font_size'])
            elif block['type'] == BlockType.IMAGE.value:
                save_base64_image(
                    block['image'],block['image_ext'], f"page_{i}_image_{image_code}",
                    os.path.join(output_dir, "images")
                )
                text = f"\n\n![image](images/page_{i}_image_{image_code}.{block['image_ext']})\n\n"
                image_code += 1
                pass
            elif block['type'] == BlockType.TABLE.value:
                text = f"\n\n{block['text']}\n\n"
            else:
                text = ""
            contents.append(text)
            pass
        # end of a page
        page_content_list.append("".join(contents))
        pass

    # 保存文件
    with open(md_path, "w", encoding='utf8') as f:
        for content in page_content_list:
            f.write(content)
            if with_page_line:
                f.write("\n\n------------\n\n")
                pass
            pass
        pass
    return md_path,md_dir


if __name__ == "__main__":
    import time

    file_path = r"C:\MyScripts\Indie\goweb\gcnote\test\docs\关于下发新春金牛购机活动的业务通知.pdf"
    output_dir = r"tmp"
    start = time.time()
    pdf_convert(file_path, output_dir, True)
    print(f"cost time: {time.time()-start}")