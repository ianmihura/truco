import csv
from data import *


def main():
    data = list()
    with open('./../truco_strength.csv', mode='r') as f:
        reader = csv.DictReader(f)
        for row in reader:
            data.append(row)

    hands = [" ".join(d['hand']) for d in data]
    scores = [d['truco_score'] for d in data]

    plt.figure(figsize=(12, 6))
    plt.plot(hands, scores)
    plt.xlabel('Hand')
    plt.ylabel('Truco Score')
    plt.xticks(rotation=90)
    plt.tight_layout()
    plt.show()


if __name__ == "__main__":
    main()
