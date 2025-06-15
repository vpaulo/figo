# Figo
Figma Golang utility lib.
This is a tool to help creating HTML and CSS from Figma designs. It does not generate clean HTML or CSS,
user input/changes are still necessary for building the desired frontend.

## Tokens
Create css variables based on figma tokens.
Figma variables can't be converted to css tokens without enterprise account to be able to use the REST API for Variables.

## CSS
Export components styles to help frontends creation, it can have repetition depending of how components
are built in Figma, use the export as a helper or a guide to accelerate development.

## HTML
Export components HTML to help frontends creation, the HTML generated gives a base structure for the
component created in Figma.
