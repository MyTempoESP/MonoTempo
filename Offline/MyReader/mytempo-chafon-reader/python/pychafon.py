import pychafon
import threading

def on_tag(epc: str, antenna: int):
    print(f"{epc}:{antenna}")

pychafon.setTagCallback(on_tag)

pychafon.read(0)

