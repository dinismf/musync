# Context Usage in Database Operations

## Why Use Context?

Context in Go is a powerful mechanism for controlling the lifecycle of operations, especially those that might take a long time to complete. In database operations, context can be used for:

1. **Cancellation**: Allows canceling operations if they take too long or if the client disconnects.
2. **Timeouts**: Sets deadlines for operations to complete.
3. **Request-scoped values**: Carries request-specific data through the call chain.
4. **Tracing and monitoring**: Helps track operations across service boundaries.

## When to Use Context vs. Without Context

### Use Context-Based Methods When:

1. **Handling HTTP Requests**: Each HTTP request should have its own context that can be canceled when the client disconnects.
   ```
   func (h *Handler) GetUser(c *gin.Context) {
       var user models.User
       if err := h.db.First(c.Request.Context(), &user, id); err != nil {
           // Handle error
       }
       // ...
   }
   ```

2. **Long-Running Operations**: Operations that might take a long time should be cancellable.
   ```
   ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
   defer cancel()
   if err := db.Find(ctx, &users, "status = ?", "active"); err != nil {
       // Handle error
   }
   ```

3. **Propagating Request-Specific Values**: When you need to pass values like user ID, trace ID, etc.
   ```
   ctx := context.WithValue(context.Background(), "user_id", userID)
   if err := db.Where(ctx, "owner_id = ?", ctx.Value("user_id")).Find(&items); err != nil {
       // Handle error
   }
   ```

### Use WithoutContext Methods When:

1. **Background Operations**: Tasks that run independently of any request.
   ```
   func backgroundTask() {
       var items []models.Item
       if err := db.FindWithoutContext(&items, "status = ?", "pending"); err != nil {
           // Handle error
       }
       // Process items...
   }
   ```

2. **Simple CRUD Operations**: For simple operations where cancellation is not critical.
   ```
   func createDefaultSettings(userID uint) {
       settings := models.Settings{UserID: userID, Theme: "default"}
       db.CreateWithoutContext(&settings)
   }
   ```

3. **Initialization Code**: Code that runs during application startup.
   ```
   func initializeDatabase() {
       var count int64
       db.DB.Model(&models.User{}).Count(&count)
       if count == 0 {
           admin := models.User{Username: "admin", Email: "admin@example.com"}
           db.CreateWithoutContext(&admin)
       }
   }
   ```

## Best Practices

1. **Always prefer context-based methods in request handlers**: This ensures that database operations can be canceled if the client disconnects.

2. **Use WithoutContext methods sparingly**: Only use them when you don't have a context available or when the operation is not tied to a specific request.

3. **Consider creating a background context for non-request operations**: Instead of using the stored context, create a background context with appropriate timeouts.
   ```
   ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
   defer cancel()
   db.Find(ctx, &items)
   ```

4. **Be consistent**: If you're using context in one part of your request handling, use it throughout the entire request lifecycle.

## Conclusion

The database package in this project provides both context-based and context-free methods to accommodate different use cases. By understanding when to use each type, you can write more robust and maintainable code.

For most HTTP request handlers, you should use the context-based methods and pass the request's context. For background tasks or initialization code, the WithoutContext methods provide a convenient alternative.
