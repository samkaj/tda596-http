# Laboration 1: HTTP server

[Lab instructions](https://chalmers.instructure.com/courses/26458/pages/lab-1-http-server)

## Prereqs

1. Install [golang](https://go.dev/doc/install). Use go 1.21.x.
1. Install [docker](https://www.docker.com/get-started).

## Running

- Run [`run.sh`](./run.sh) if you have bash.
- If no bash: [`run.bat`](./run.bat) in powershell on windows.

## Requirements

- Allowed values of `Content-Type`
  - text/html, text/plain, image/gif, image/jpeg, image/jpeg, or text/css.
  - If not in above list, return `400 - Bad request`.
- Allow `GET`:
  - Other methods result in `501 - Not implemented`.
  - `GET` on non-existent files results in `404 - Not found`.
  - Other bad requests, such as badly formatted headers, return appropriate `4XX` messages.
- Allow `POST`:
  - Store files appropriately and make the accessible via subsequent `GET` requests.
  - All files should be stored in `/fs`
- Use the `http` package
  - Do not use `ListenAndServe` or `Listen` followed by `Serve`.
- Concurrency:
  - Allow at most 10 concurrent clients.
