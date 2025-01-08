import json
import pymupdf
from pymupdf import Document
from pymupdf import Page
from pdf_parser.block_types import BlockType


def read_as_dict(file_path: str) -> list[dict]:
    page_content = []
    doc = pymupdf.open(file_path)
    for i in range(0, doc.page_count):
        page = doc[i]
        fitz_page_dict = json.loads(page.get_text("json", sort=True))
        page_dict = {
            "page_number": i,
            "chunks": []
        }
        pass
        for block in fitz_page_dict['blocks']:
            page_block = {
                "font_size": 24,  # h1是48 = 24 *2 超过了也是h1无所谓了
                "bbox": block['bbox'],
            }
            # 文字
            if 'lines' in block:
                page_block['type'] = BlockType.TEXT.value
                page_block_text_list = []
                for line in block['lines']:
                    for span in line['spans']:
                        page_block_text_list.append(span["text"])
                        if span['size'] < page_block['font_size']:
                            page_block['font_size'] = span['size']
                # 回到block级别
                page_block["text"] = "".join(page_block_text_list)
                pass
            elif 'image' in block:
                page_block['type'] = BlockType.IMAGE.value
                page_block['image'] = block['image']
                page_block['image_ext'] = block['ext']
                pass

            page_dict["chunks"].append(page_block)
        pass
        # 一个页面接受
        page_content.append(page_dict)
    #
    return page_content


if __name__ == "__main__":
    import time

    file_path = r"E:\00-Document\Indie\GoWeb\pdf2md\example\23年统计公报-节选.pdf"
    start = time.time()
    read_as_dict(file_path)
    print(time.time() - start)
