import os
import re
from io import StringIO, BytesIO
import base64
from PIL import Image


def save_base64_image(base64_string: str, ext: str, filename: str, output_dir: str) -> str:
    """
    将 base64 编码的图片字符串保存为指定扩展名的图片文件。

    :param base64_string: 包含图片数据的 base64 字符串（不包含 data URL scheme 前缀）
    :param ext: 图片文件的扩展名，例如 'png', 'jpg' 等
    :param filename: 保存到磁盘的文件名（不含路径和扩展名）
    :param output_dir: 输出目录路径
    :return: 保存成功的文件路径
    :raises ValueError: 如果 base64 解码失败或图像无法保存
    """
    try:
        # 解码 base64 字符串
        image_data = base64.b64decode(base64_string)

        # 创建一个 BytesIO 对象
        image_io = BytesIO(image_data)

        # 打开图像
        image = Image.open(image_io)

        # 确保输出目录存在
        if not os.path.exists(output_dir):
            os.makedirs(output_dir)

        # 保存图像
        output_filename = f"{filename}.{ext}"
        output_path = os.path.join(output_dir, output_filename)
        image.save(output_path)

        # 验证文件是否存在
        if not os.path.exists(output_path):
            raise IOError("文件保存失败，路径无效或其他问题")

        return output_path

    except Exception as e:
        raise ValueError(f"保存图像时出错: {e}")


def get_page_image(image: Image, block: dict, dpi: int):
    try:
        r_box = tuple([b * dpi / 72 for b in block['bbox']]) # 默认是72
        image_clip = image.crop(r_box)
        # 将图像保存到内存中的 BytesIO 对象
        buffered = BytesIO()
        image_clip.save(buffered, format="PNG")
        # 获取图像的 Base64 编码
        encoded_image = base64.b64encode(buffered.getvalue()).decode("utf-8")
        return encoded_image
    except Exception as e:
        return block['image']


def word2heading(text: str, font_size: float):
    # 非常开心，pdf的字号标准也是pt，那么就参考docx2md的方式了
    # 参考 https://github.com/zakahan/docx2md word2Heading部分
    font_size = font_size * 2
    maxHeadingLength = 15
    h1 = 48
    h2 = 36
    h3 = 28
    h4 = 24
    if check_string_prefix(text):
        if h1 <= font_size:
            return "\n\n# " + text + "\n\n"
        elif h2 <= font_size:
            return "\n\n## " + text + "\n\n"
        elif h3 <= font_size:
            return "\n\n### " + text + "\n\n"
        elif h4 <= font_size and len(text) < maxHeadingLength:
            return "\n\n#### " + text + "\n\n"
        else:
            return text
    else:
        if h1 <= font_size:
            return "\n\n# " + text + "\n\n"
        elif h2 <= font_size:
            return "\n\n## " + text + "\n\n"
        elif h3 <= font_size and len(text) < maxHeadingLength:
            return "\n\n### " + text + "\n\n"
        else:
            return text


def check_string_prefix(s: str) -> bool:
    """
    检查字符串是否满足以下条件之一：
    1. 是否以汉字数字开头
    2. 是否以阿拉伯数字开头
    3. 是否以 "第" + 阿拉伯数字开头
    4. 是否以 "第" + 汉字数字开头

    :param s: 输入字符串
    :return: 如果满足任意条件返回 True，否则返回 False
    """
    # 汉字数字的正则表达式
    chinese_digits = "[一二三四五六七八九十百千万亿]"

    # 条件 1: 汉字数字开头
    if re.match(f"^{chinese_digits}", s):
        return True

    # 条件 2: 阿拉伯数字开头
    if re.match(r"^\d", s):
        return True

    # 条件 3: "第" + 阿拉伯数字开头
    if re.match(r"^第\d", s):
        return True

    # 条件 4: "第" + 汉字数字开头
    if re.match(f"^第{chinese_digits}", s):
        return True

    return False


if __name__ == "__main__":
    x = check_string_prefix("第个")
    print(x)