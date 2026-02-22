# WebMCP Tool Patterns

This reference provides detailed patterns for implementing WebMCP tools in web applications.

## Tool Registration API

### Basic Registration

```javascript
// Register a single tool
const unregister = navigator.webMCP.registerTool({
  name: "searchProducts",
  description: "Search the product catalog by name, category, or attributes",
  inputSchema: {
    type: "object",
    properties: {
      query: {
        type: "string",
        description: "Search query string"
      },
      category: {
        type: "string",
        description: "Product category filter"
      },
      inStock: {
        type: "boolean",
        description: "Filter to only show in-stock items"
      }
    },
    required: ["query"]
  }
}, async (params) => {
  const results = await productSearch(params);
  updateProductDisplay(results);
  return { products: results, count: results.length };
});

// Unregister when no longer needed
unregister();
```

### Multiple Tools

```javascript
// Register multiple related tools
const tools = [
  {
    name: "getProduct",
    description: "Get details for a specific product",
    inputSchema: {
      type: "object",
      properties: {
        productId: { type: "string" }
      },
      required: ["productId"]
    }
  },
  {
    name: "addToCart",
    description: "Add a product to the shopping cart",
    inputSchema: {
      type: "object",
      properties: {
        productId: { type: "string" },
        quantity: { type: "integer", minimum: 1, default: 1 }
      },
      required: ["productId"]
    }
  }
];

// Register all tools
const unregisterAll = navigator.webMCP.registerTools(tools, handlers);
```

## Parameter Schema Patterns

### String with Enum

```javascript
{
  color: {
    type: "string",
    enum: ["red", "blue", "green", "yellow", "black", "white"],
    description: "Color filter for products"
  }
}
```

### Numeric with Range

```javascript
{
  priceRange: {
    type: "object",
    properties: {
      min: { type: "number", minimum: 0 },
      max: { type: "number", maximum: 10000 }
    }
  }
}
```

### Array of Items

```javascript
{
  tags: {
    type: "array",
    items: { type: "string" },
    description: "List of tags to filter by"
  }
}
```

### Nested Object

```javascript
{
  address: {
    type: "object",
    properties: {
      street: { type: "string" },
      city: { type: "string" },
      zipCode: { type: "string", pattern: "^\\d{5}(-\\d{4})?$" }
    },
    required: ["street", "city", "zipCode"]
  }
}
```

## Handler Implementation Patterns

### Using Existing Page Logic

```javascript
// The handler should leverage existing application code
async function handleSearchProducts(params) {
  // Reuse existing search function from the page
  const searchFn = window.__app.searchProducts;
  if (!searchFn) {
    throw new Error("Search functionality not available");
  }

  // Call with validated parameters
  const results = await searchFn({
    query: params.query,
    filters: {
      category: params.category,
      inStock: params.inStock
    }
  });

  // Update UI so user sees the results
  window.__app.updateProductList(results);

  // Return structured data for the agent
  return {
    products: results.map(p => ({
      id: p.id,
      name: p.name,
      price: p.price,
      inStock: p.stock > 0
    })),
    totalCount: results.length,
    pageUrl: window.location.href
  };
}
```

### Async Operations with Progress

```javascript
async function handleGenerateReport(params) {
  // Start the operation
  const jobId = await reportService.startGeneration(params);

  // Poll for completion
  let attempts = 0;
  while (attempts < 30) {
    const status = await reportService.getStatus(jobId);

    if (status.complete) {
      // Show result to user
      displayReport(status.result);

      return {
        success: true,
        reportUrl: status.downloadUrl,
        generatedAt: status.completedAt
      };
    }

    if (status.failed) {
      return {
        success: false,
        error: status.error
      };
    }

    // Wait before next poll
    await new Promise(r => setTimeout(r, 1000));
    attempts++;
  }

  return {
    success: false,
    error: "Report generation timed out"
  };
}
```

### Human Confirmation Pattern

```javascript
async function handleDeleteAccount(params) {
  // Show confirmation dialog in the page
  const confirmed = await showConfirmDialog({
    title: "Delete Account?",
    message: "This action cannot be undone. All your data will be permanently deleted.",
    confirmLabel: "Delete My Account",
    cancelLabel: "Cancel",
    destructive: true
  });

  if (!confirmed) {
    return {
      success: false,
      reason: "User cancelled the operation"
    };
  }

  // Proceed with deletion
  await accountService.delete(params.accountId);

  // Update UI
  navigateTo("/goodbye");

  return {
    success: true,
    message: "Account deleted successfully"
  };
}
```

### Form Input Collection

```javascript
async function handleScheduleAppointment(params) {
  // Show form to collect missing information
  const formData = await showFormDialog({
    title: "Schedule Appointment",
    fields: [
      {
        name: "date",
        type: "date",
        label: "Preferred Date",
        required: true,
        min: tomorrow()
      },
      {
        name: "time",
        type: "select",
        label: "Time Slot",
        options: ["9:00 AM", "10:00 AM", "2:00 PM", "3:00 PM"],
        required: true
      },
      {
        name: "notes",
        type: "textarea",
        label: "Additional Notes"
      }
    ]
  });

  if (!formData) {
    return { success: false, reason: "User cancelled" };
  }

  // Book the appointment
  const appointment = await bookingService.create({
    ...params,
    ...formData
  });

  // Update calendar UI
  calendarView.addEvent(appointment);

  return {
    success: true,
    appointment: {
      id: appointment.id,
      date: appointment.date,
      time: appointment.time
    }
  };
}
```

## Error Handling Patterns

### Input Validation

```javascript
async function handleTransfer(params) {
  // Validate required parameters
  if (!params.recipientId) {
    return {
      success: false,
      error: "Recipient is required",
      errorCode: "MISSING_RECIPIENT"
    };
  }

  // Validate constraints
  if (params.amount <= 0) {
    return {
      success: false,
      error: "Amount must be greater than zero",
      errorCode: "INVALID_AMOUNT"
    };
  }

  // Validate business rules
  const balance = await getAccountBalance();
  if (params.amount > balance) {
    return {
      success: false,
      error: `Insufficient funds. Current balance: ${formatCurrency(balance)}`,
      errorCode: "INSUFFICIENT_FUNDS"
    };
  }

  // Proceed with transfer
  // ...
}
```

### Graceful Degradation

```javascript
async function handleGetData(params) {
  try {
    const data = await fetchData(params);
    return { success: true, data };
  } catch (networkError) {
    // Try fallback data source
    try {
      const cachedData = await getCachedData(params);
      return {
        success: true,
        data: cachedData,
        warning: "Using cached data. Some information may be outdated."
      };
    } catch (cacheError) {
      return {
        success: false,
        error: "Unable to retrieve data. Please try again later.",
        canRetry: true
      };
    }
  }
}
```

## Integration with Page Lifecycle

### Tool Availability Management

```javascript
// Tools that are only available on certain pages
class ToolManager {
  constructor() {
    this.registeredTools = new Map();
  }

  enableTool(name, definition, handler) {
    if (this.registeredTools.has(name)) {
      return; // Already registered
    }

    const unregister = navigator.webMCP.registerTool(definition, handler);
    this.registeredTools.set(name, unregister);
  }

  disableTool(name) {
    const unregister = this.registeredTools.get(name);
    if (unregister) {
      unregister();
      this.registeredTools.delete(name);
    }
  }

  disableAll() {
    for (const unregister of this.registeredTools.values()) {
      unregister();
    }
    this.registeredTools.clear();
  }
}

// Usage based on page state
const toolManager = new ToolManager();

// Enable editing tools when document is loaded
documentEditor.onDocumentLoad((doc) => {
  toolManager.enableTool("editDocument", editDocumentDef, handleEditDocument);
  toolManager.enableTool("saveDocument", saveDocumentDef, handleSaveDocument);
});

// Disable when document is closed
documentEditor.onDocumentClose(() => {
  toolManager.disableTool("editDocument");
  toolManager.disableTool("saveDocument");
});
```

### React Integration

```jsx
import { useEffect } from 'react';

function useWebMCPTool(name, definition, handler, deps = []) {
  useEffect(() => {
    if (!navigator.webMCP) {
      console.warn('WebMCP not available');
      return;
    }

    const unregister = navigator.webMCP.registerTool(
      { name, ...definition },
      handler
    );

    return unregister;
  }, [name, ...deps]);
}

// Usage in component
function ProductPage({ productId }) {
  useWebMCPTool(
    "addReview",
    {
      description: "Add a review for this product",
      inputSchema: {
        type: "object",
        properties: {
          rating: { type: "integer", minimum: 1, maximum: 5 },
          comment: { type: "string" }
        },
        required: ["rating"]
      }
    },
    async (params) => {
      const review = await submitReview(productId, params);
      return { success: true, review };
    },
    [productId]
  );

  return <ProductDisplay productId={productId} />;
}
```

## Security Best Practices

### Sanitize Agent Input

```javascript
async function handleSearch(params) {
  // Sanitize string inputs
  const sanitizedQuery = sanitizeString(params.query);

  // Validate and coerce types
  const limit = Math.min(Math.max(parseInt(params.limit) || 10, 1), 100);

  // Use parameterized queries
  const results = await db.query(
    "SELECT * FROM products WHERE name LIKE ? LIMIT ?",
    [`%${sanitizedQuery}%`, limit]
  );

  return { results };
}
```

### Rate Limiting

```javascript
const rateLimiter = new Map();

async function withRateLimit(userId, action, handler) {
  const key = `${userId}:${action}`;
  const now = Date.now();
  const lastCall = rateLimiter.get(key) || 0;

  if (now - lastCall < 1000) { // 1 second minimum between calls
    return {
      success: false,
      error: "Please wait before trying again",
      errorCode: "RATE_LIMITED"
    };
  }

  rateLimiter.set(key, now);
  return handler();
}
```

### Audit Logging

```javascript
async function handleSensitiveAction(params, context) {
  // Log the attempt
  await auditLog.record({
    action: "sensitive_action",
    userId: context.userId,
    params: sanitizeForLog(params),
    timestamp: new Date().toISOString(),
    userAgent: context.agentId
  });

  // Execute the action
  const result = await executeSensitiveAction(params);

  // Log the result
  await auditLog.record({
    action: "sensitive_action_complete",
    userId: context.userId,
    success: result.success,
    timestamp: new Date().toISOString()
  });

  return result;
}
```
