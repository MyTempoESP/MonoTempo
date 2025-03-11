from mytempo_reader.mytempo_reader import Reader
from rabbit.rabbit import TagSender

if __name__ == "__main__":
    with TagSender() as ts:
        reader = Reader(ts)
        reader.start()
