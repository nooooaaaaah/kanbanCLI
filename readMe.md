# Kanban Board CLI

This project is a command-line interface (CLI) for managing a Kanban board. It's written in Go and uses a JSON file for data persistence.

## Features

- Create, update, and delete cards and lists.
- Navigate through the board using keyboard shortcuts.
- Persist data in a JSON file.

## Installation

Ensure you have Go installed on your machine. You can download it from the [official Go website](https://golang.org/dl/).

Clone the repository:

```sh
git clone https://github.com/yourusername/kanban-cli.git
```

Navigate to the project directory:

```sh
cd kanban-cli
```

Build the project:

```sh
go build
```

## Usage

Run the application:

```sh
./kanban-cli
```

Here are the keyboard shortcuts you can use to navigate through the board:

- `left` or `h`: Move the cursor to the left.
- `right` or `l`: Move the cursor to the right.
- `i`: Toggle the input field for adding a new card.
- `enter`: Toggle the selection of the currently highlighted list.
- `ctrl+c` or `q`: Quit the application.

When the input field for adding a new card is displayed, type the title of the new card and press [`/`](command:_github.copilot.openRelativePath?%5B%22%2F%22%5D "/") to add the card.

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

## License

[MIT](https://choosealicense.com/licenses/mit/)

