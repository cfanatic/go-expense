# go-expensegui

This program is the graphical front-end for [go-expense](https://github.com/cfanatic/go-expense).

It helps you analyze your monthly and yearly expenses by giving answers to following questions:

1. How much do you spend on average per month and year?
2. How much would you be able to set aside each month for savings?
3. What are the highest cost types, e.g. rent, car, food, etc.?
4. Where are potentials to reduce costs?

The graphical user interface shall simplify the process of importing and labeling account statement transactions.
Labeling is done in automatic fashion whenever a similar transaction can be found in the database.

Three sample account statements can be found in folder `res/`.

## Requirements

Developed and tested on the following setup:

- macOS 10.15.2
- Go 1.13.4
- Docker 2.2.0.0

## Installation

Run the particular build process for one of the hosts below:

### macOS

```bash
qtmoc desktop
go build
```

### Linux

```bash
docker pull therecipe/qt:linux_debian_9
docker build -t expensegui:latest -f Dockerfile .
docker run --name expensegui expensegui:latest
docker cp expensegui:/home/user/work/src/github.com/cfanatic/go-expensegui/deploy/linux/go-expensegui .
```

### Windows

```bash
docker pull therecipe/qt:windows_64_static
docker build -t expensegui:latest -f Dockerfile .
docker run --name expensegui expensegui:latest
docker cp expensegui:/home/user/work/src/github.com/cfanatic/go-expensegui/deploy/windows/go-expensegui.exe .
```

## Usage

Start the MongoDB database:

```bash
docker pull mongo:latest
docker run -d -p 27017:27017 --name expenses mongo:latest
```

Run the program:

```bash
./go-expensegui
```

Import each account statement and specify cost type labels for all transactions:

![Data_Label](https://raw.githubusercontent.com/cfanatic/go-expensegui/master/res/go-expensegui-1.png)

Show transaction details by hovering over Payee fields:

![Data_Exclude](https://raw.githubusercontent.com/cfanatic/go-expensegui/master/res/go-expensegui-3.png)

Exclude cost groups from statistical analysis by right-clicking on Label fields:

![Data_Exclude](https://raw.githubusercontent.com/cfanatic/go-expensegui/master/res/go-expensegui-2.png)

## Analysis

Analyze the monthly expenses for a particular account statement:

![Data_Exclude](https://raw.githubusercontent.com/cfanatic/go-expensegui/master/res/go-expensegui-4.png)

Show the complete analysis for all account statements which are available in the database:

![Data_Exclude](https://raw.githubusercontent.com/cfanatic/go-expensegui/master/res/go-expensegui-5.png)
