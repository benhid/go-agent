import sys
from handler import handler

if __name__ == "__main__":
    stdin = sys.stdin.readlines()
    handler.handle(" ".join(stdin))
