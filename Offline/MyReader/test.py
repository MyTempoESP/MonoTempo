from rabbit.rabbit import Checkup

check = Checkup()


def handler(ch, method, properties, body):
    print(body)


check.handler = handler


def main():
    with check as c:
        try:
            c.start()
        except KeyboardInterrupt:
            return


if __name__ == "__main__":
    main()
