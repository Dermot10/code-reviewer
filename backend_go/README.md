Design Decisions -

# Multiple operations for different transport / interaction styles -

## WS for live updates - collaborative file updates

## HTTP for static file updates using standard CRUD over a standard HTTP connection to interact with resources

## Both utilise the same service level business logic.

### They differ at infrastruture level as they use different client level implementations. This is common for slack, notion, figma etc
