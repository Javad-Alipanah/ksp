# Knight Shortest Path


###### This is the solution to the knight's shortest path problem in Golang

# Table of Contents

1. [Problem Description](#problem-description)
2. [How to Run](#how-to-run)
3. [GraphQL Endpoint](#graphql-endpoint)
4. [Example](#example)

## Problem Description
The problem is as follows:

    Assume that you have an N*N chess board with a knight in some cell called "start" point, and an arbitrary cell called "target" point.
    Calculate the shortest path for the knight from start to target (the knight may move two cells vertically and one cell horizontally, or two cells horizontally and one cell vertically).

## How to Run
##### Prerequisites: `docker`, `docker-compose`, `build-essential`, `go`
1. Clone the project:
    ```bash
    git clone git@github.com:Javad-Alipanah/ksp.git && cd ksp
    ```
2. Run:
    ```bash
    sudo docker-compose up -d
    ```
Now the GraphQL endpoint is available at [localhost](http://localhot)
 
To see the project logs run:
```bash
sudo docker-compose logs [-f]
```

# GraphQL Endpoint
The project has a GraphQL endpoint by default at [localhost](http://localhot) from which you can interact with the application and calculate the answer for different outputs.

If you are not familiar with GraphQL you can visit [here](https://graphql.org/) for introductory info.

If you visit the project root at [localhost](http://localhot) you can see an interactive console from which you can execute these queries and mutations:
###### Note: The (0, 0) cell is the upper left cell in your board

## Queries
### Boards:
```graphql
query {
  boards {
      # subfield selection
  }
}
```
Full Example:
```graphql
query {
  boards {
    id
    size
    start {
      x
      y
    }
    target {
      x
      y
    }
    path {
      x
      y
    }
  }
}
```

### Board
```graphql
query {
  board(id: $someId) {
    # subfield selection
  }
}
```

Full Example:
```graphql
query {
  board(id: 4) {
    id
    size
    start {
      x
      y
    }
    target {
      x
      y
    }
    path {
      x
      y
    }
  }
}
```

## Mutations

### Create Board
###### Note: The path calcultaion is asynchronous; so when you create or update a board, the returned path is always empty. You must qury the board to get the calculated shortest path. If the path is still empty then no route from "start" to "target" exists.
```graphql
mutation {
  createBoard(input:{size: $size, start: {x: $x1, y: $y1}, target: {x: $x2, y: $y2}}) {
    # subfield selection
  }
}
```
Full Example:
```graphql
mutation {
  createBoard(input:{size: 4, start: {x: 1, y:1}, target: {x: 2, y:2}}) {
    id
    size
    start {
      x
      y
    }
    target {
      x
      y
    }
    path {
      x
      y
    }
  }
}
```
### Update Board
```graphql
mutation {
  updateBoard(board: $partialBoardJsonWithIdProvided) {
    # subfield selection
  }
}
```

Full Example:
```graphql
mutation {
  updateBoard(board:{id: 6, size: 6}) {
    id
    size
    start {
      x
      y
    }
    target {
      x
      y
    }
    path {
      x
      y
    }
  }
}
```

### Delete Board
###### Note: No subfield selection is allowed. The only returned data is either an error or a boolean indicating whether the delete succeded or not.
```graphql
mutation {
  deleteBoard(id: $someId)
}
```
Example:
```graphql
mutation {
  deleteBoard(id: 6)
}
```

# Example:
1. Create the Board:
    ```graphql
    mutation {
      createBoard(
        board: {
          size: 15,
          start: { x: 3, y: 7 },
          target: { x: 6, y: 0 }
        }
      ) 
      {
        id
      }
    }
    ```

2. Query the Board:
    ```graphql
    query {
      board(id: $idFromPrevStep) {
        path {
          x
          y
        }
      }
    }
    ```

3. Verify the Output:
    ```json
    {
      "data": {
        "board": {
          "path": [
            {
              "x": 3,
              "y": 7
            },
            {
              "x": 4,
              "y": 5
            },
            {
              "x": 5,
              "y": 3
            },
            {
              "x": 4,
              "y": 1
            },
            {
              "x": 6,
              "y": 0
            }
          ]
        }
      }
    }
    ```