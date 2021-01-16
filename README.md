# ShakeSearch

This is simple web app that allows a user to search for a text string in
the complete works of Shakespeare.

## Try it!
https://rubbioli.com/shakesearch

## Note
This is a fork of https://github.com/ProlificLabs/shakesearch.

## Changes

- Project structure to make it more readable and testable
- Search is now case-insensitive
- Added a cache layer to make most common searches faster
- Added some unit tests to make it more reliable
- Fixed `slice out of bounds` error sometimes found when searching

## Next steps

- Add fuzzy search
