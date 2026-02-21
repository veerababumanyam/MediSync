---
name: webmcp
description: This skill should be used when the user asks to "expose web app functionality to AI agents", "create client-side MCP tools", "implement WebMCP", "register JavaScript tools for AI", "human-in-the-loop web automation", "browser agent integration", or mentions WebMCP tools, MCP-B, or exposing site capabilities as tools.
---

# WebMCP - Web Model Context Protocol

WebMCP is a proposed web API that enables web applications to expose their functionality as JavaScript-based tools accessible to AI agents and assistive technologies. It brings the Model Context Protocol (MCP) to the client side, enabling human-in-the-loop collaborative workflows.

★ Insight ─────────────────────────────────────
WebMCP's key innovation:
1. **Client-side tools** - JavaScript functions as AI-callable tools
2. **Shared context** - User, page, and agent operate on same state
3. **HITL by design** - Humans maintain visibility and control

This allows websites to become "MCP servers" without backend changes.
─────────────────────────────────────────────────

## Quick Reference

| Aspect | Details |
|--------|---------|
| **Purpose** | Expose web app functionality to AI agents |
| **Status** | Proposed standard (explainer stage) |
| **Authors** | Microsoft, Google |
| **Key Benefit** | Reuse frontend code for AI tool calls |
| **Protocol** | Compatible with MCP |

## Core Concept

WebMCP allows web pages to register tools that AI agents can discover and invoke:

```javascript
// Web page registers tools
navigator.webMCP.registerTool({
  name: "filterProducts",
  description: "Filter product list based on criteria",
  parameters: {
    type: "object",
    properties: {
      category: { type: "string" },
      maxPrice: { type: "number" }
    }
  }
}, async (params) => {
  // Tool implementation uses existing page logic
  const results = await productStore.filter(params);
  updateUI(results);
  return results;
});
```

## Why WebMCP?

### Advantages Over Backend MCP

| Aspect | Backend MCP | WebMCP |
|--------|-------------|--------|
| **Implementation** | New server code | Reuse frontend JS |
| **Context** | Separate from UI | Shared with user |
| **Auth** | New auth layer | Use existing session |
| **State** | Sync required | Already in sync |
| **Visibility** | Hidden operations | User can observe |

### Human-in-the-Loop Workflows

Unlike fully autonomous agents, WebMCP enables:
- User starts a task, agent helps complete it
- Agent takes action, user reviews and approves
- User and agent collaborate in real-time
- All actions visible in the shared UI

## Tool Definition Structure

```javascript
/**
 * Tool definition with JSDoc for natural language description
 *
 * @description Filter products based on user criteria
 * @param {string} category - Product category to filter by
 * @param {number} maxPrice - Maximum price in USD
 * @returns {Product[]} Filtered product list
 */
function filterProducts(category, maxPrice) {
  // Implementation uses existing page logic
}
```

### Schema Definition

```javascript
const tool = {
  name: "getDresses",
  description: "Returns product listings with id, description, price, photo",
  inputSchema: {
    type: "object",
    properties: {
      size: {
        type: "number",
        description: "EU dress size between 2 and 14"
      },
      color: {
        type: "string",
        enum: ["Red", "Blue", "Green", "Yellow", "Black", "White"]
      }
    }
  }
};
```

## Use Cases

### 1. Shopping Assistance

User browses an e-commerce site and asks agent for help:

```
User: Show me cocktail dresses in size 6 under $200

Agent calls: getDresses({ size: 6, maxPrice: 200, style: "cocktail" })
UI updates: Shows filtered results
Agent: "Here are 5 cocktail dresses in your size under $200."
```

### 2. Creative Applications

User works on a design tool:

```
User: Make the heading red and add a spring-themed background

Agent calls: editDesign("Change heading color to red, add spring background")
UI updates: Design changes in real-time
User: Reviews and adjusts
```

### 3. Code Review

User reviews a pull request:

```
User: Why is the Mac bot failing?

Agent calls: getTryRunStatuses()
Agent calls: getTryRunFailureSnippet({ bot_name: "mac-x64-rel" })
Agent: "The Mac bot is failing with an 'Out of Space' error..."
```

## Architecture Comparison

### Backend MCP (Traditional)
```
Agent → MCP Server → Backend API → Database
         (separate infrastructure)
```

### WebMCP
```
Agent → Browser Tab → WebMCP Tools → Existing Page Logic → UI
         (shared context, visible to user)
```

## Integration with MCP

WebMCP is designed to work alongside MCP:

```javascript
// WebMCP tools can be exposed as MCP tools
const webMCPServer = {
  tools: [/* client-side tools */],
  resources: [/* page resources */],
  prompts: [/* helpful prompts */]
};

// Browser can expose these to external MCP clients
navigator.mcp.expose(webMCPServer);
```

## Security Considerations

### Permissions Model

1. **Site Registration** - User grants site permission to register tools
2. **Agent Access** - User grants agent permission to use tools
3. **Per-call Visibility** - User sees what data is sent/received

### Trust Boundaries

```javascript
// Browser prompts user when:
// 1. Site first registers tools
// 2. New agent requests tool access
// 3. Tool receives/sends sensitive data

// User can:
// - Allow/deny per tool
// - Set always-allow for trusted pairs
// - Review tool call history
```

### Cross-Origin Isolation

Data from one origin's tools may flow to another origin's tools. Browser should:
- Show which apps are being invoked
- Display data being transferred
- Allow user intervention

## Implementing WebMCP Tools

### Basic Pattern

```javascript
// 1. Define tool with schema
const myTool = {
  name: "searchInventory",
  description: "Search inventory for items matching criteria",
  inputSchema: {
    type: "object",
    properties: {
      query: { type: "string" },
      inStock: { type: "boolean" }
    },
    required: ["query"]
  }
};

// 2. Implement handler
async function handleSearch(params) {
  // Use existing page logic
  const results = await inventoryAPI.search(params.query);

  // Update UI (visible to user)
  renderSearchResults(results);

  // Return structured data
  return {
    products: results.map(r => ({
      id: r.id,
      name: r.name,
      price: r.price,
      stock: r.stock
    }))
  };
}

// 3. Register with WebMCP
navigator.webMCP.registerTool(myTool, handleSearch);
```

### HITL Pattern

```javascript
const deleteTool = {
  name: "deleteDocument",
  description: "Delete a document after user confirmation",
  inputSchema: {
    type: "object",
    properties: {
      documentId: { type: "string" }
    },
    required: ["documentId"]
  }
};

async function handleDelete({ documentId }) {
  // Show confirmation UI
  const confirmed = await showConfirmDialog({
    title: "Delete Document?",
    message: "This action cannot be undone.",
    options: ["Delete", "Cancel"]
  });

  if (!confirmed) {
    return { success: false, reason: "User cancelled" };
  }

  // Proceed with deletion
  await documentAPI.delete(documentId);
  updateUI();

  return { success: true };
}
```

## Relationship to Other Protocols

| Protocol | Purpose | WebMCP Relationship |
|----------|---------|---------------------|
| **MCP** | Backend tool integration | WebMCP brings MCP concepts to client |
| **A2A** | Agent-to-agent communication | WebMCP is agent-to-web-page |
| **OpenAPI** | HTTP API description | WebMCP tools are JavaScript, not HTTP |

## Future Directions

### Progressive Web Apps (PWA)
- Tools declared in manifest
- Available "offline" (without page open)
- System can launch PWA for tool calls

### Background Model Context
- Tools handled in service worker
- No browser window required
- Notification-based feedback

## Best Practices

1. **Reuse existing code** - Don't duplicate, leverage page logic
2. **Keep tools focused** - Single responsibility per tool
3. **Provide clear descriptions** - Help AI understand tool purpose
4. **Validate inputs** - Sanitize agent-provided parameters
5. **Show feedback** - Keep user informed of actions
6. **Handle errors gracefully** - Return meaningful error messages

## Key Resources

| Resource | URL |
|----------|-----|
| Explainer | https://github.com/webmachinelearning/webmcp |
| Proposal | https://github.com/webmachinelearning/webmcp/blob/main/proposal.md |
| MCP Specification | https://modelcontextprotocol.io |

## Related Skills

- **copilotkit-generative-ui** - For generative UI with CopilotKit
- **medisync-dev** - For MediSync-specific development
- **tally-integration** - For Tally ERP integration patterns
