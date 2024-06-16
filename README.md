# UI.WEB

A platform for creating web-based User Interfaces that are highly programmable.

Project Status: IN EARLY DEVELOPMENT

## Objective

UI.WEB aims to employ the [emacs](https://www.gnu.org/software/emacs/) architecture and ideas
on top of web technologies.

It does not attempt to look like a desktop application as it can be opened on any browser the user may choose.
However, by making it a PWA, it may look exactly like one.

## Project Structure

* `common/` - Source code used in both backend and frontend.
* `frontend/` - Frontend source code and assets.
* `backend/` - Backend source code.

## Building

Building this project requires node.js.

The build uses [esbuild](https://esbuild.github.io/) to bundle both the backend and the frontend.

To build:

```shell
cd backend
npm install
npm run build
```

The backend build builds the frontend first, so the above builds everything required to run UI.WEB.

To run (from the `backend` dir):

```shell
npm start
```
