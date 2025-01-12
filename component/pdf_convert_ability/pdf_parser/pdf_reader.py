import json
import pymupdf
from pymupdf import Document
from pymupdf import Page
from pymupdf.table import find_tables
from pdf_parser.block_types import BlockType


def read_as_dict(file_path: str) -> list[dict]:
    page_content = []
    doc = pymupdf.open(file_path)
    for i in range(0, doc.page_count):
        page = doc[i]
        tables = find_tables(page).tables
        table_block_list = []
        for table in tables:
            table_block_list.append(
                {
                    "font_size": 12,  # 默认常规字号
                    "bbox": list(table.bbox),
                    "text": table.to_markdown(),
                    "type": BlockType.TABLE.value
                }
            )
            # 这里开始提取表格
            # 根据box范围来确定插入位置
            pass
        fitz_page_dict = json.loads(page.get_text("json", sort=True))
        page_dict = {
            "page_number": i,
            "chunks": []
        }


        for block in fitz_page_dict['blocks']:
            page_block = {
                "font_size": 24,  # h1是48 = 24 *2 超过了也是h1无所谓了
                "bbox": block['bbox'],
            }
            # 有表格的情况
            for table_block in table_block_list:
                if is_inside(table_block["bbox"], page_block["bbox"]):      # B 是否属于 A
                    # 如果有重叠，直接不要page_block了
                    if len(page_dict["chunks"]) == 0 or page_dict["chunks"][-1] not in table_block_list:
                        page_dict["chunks"].append(table_block)
                        pass
                    else:
                        pass
                else:
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
                    else:
                        # 啥也不干，但这个情况不太可能
                        page_block["text"] = ""
                        pass

                    page_dict["chunks"].append(page_block)
                    # end of else
                # end of for loop
                pass
            # 没有表格的情况
            if len(table_block_list) == 0:
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
                else:
                    # 啥也不干，但这个情况不太可能
                    page_block["text"] = ""
                    pass

                page_dict["chunks"].append(page_block)
                pass        # end of if
            pass # end of everything
        pass

        # 一个页面接受
        page_content.append(page_dict)
    #
    return page_content


def is_intersect(rectA, rectB):
    # rectA and rectB are lists in the form [x1, y1, x2, y2]
    # where (x1, y1) is the bottom-left corner and (x2, y2) is the top-right corner of the rectangle

    A_x1, A_y1, A_x2, A_y2 = rectA
    B_x1, B_y1, B_x2, B_y2 = rectB

    # Check if one rectangle is to the left of the other
    if A_x2 < B_x1 or B_x2 < A_x1:
        return False

    # Check if one rectangle is above the other
    if A_y2 < B_y1 or B_y2 < A_y1:
        return False

    # If neither of the above, the rectangles intersect
    return True


def is_inside(rectA, rectB):
    # rectA and rectB are lists in the form [x1, y1, x2, y2]
    # where (x1, y1) is the bottom-left corner and (x2, y2) is the top-right corner of the rectangle

    A_x1, A_y1, A_x2, A_y2 = rectA
    B_x1, B_y1, B_x2, B_y2 = rectB

    # Check if all corners of rectB are inside rectA
    return (A_x1 <= B_x1 and B_x2 <= A_x2) and (A_y1 <= B_y1 and B_y2 <= A_y2)


if __name__ == "__main__":
    import time

    file_path = r"C:\MyScripts\Indie\goweb\gcnote\test\docs\23年统计公报-节选.pdf"
    start = time.time()
    x = read_as_dict(file_path)
    print(time.time() - start)
