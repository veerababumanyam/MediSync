// MediSync WebMCP Tools Example
// This example shows how to implement WebMCP tools for MediSync

// ============================================
// TOOL REGISTRATION SETUP
// ============================================

class MediSyncWebMCPManager {
  constructor() {
    this.registeredTools = new Map();
    this.sessionContext = null;
  }

  // Initialize WebMCP for the page
  async init(context) {
    this.sessionContext = context;

    // Check if WebMCP is available
    if (!navigator.webMCP) {
      console.warn('WebMCP API not available in this browser');
      return false;
    }

    // Register core tools
    this.registerCoreTools();

    // Listen for page changes to update available tools
    this.setupPageLifecycleHandlers();

    return true;
  }

  registerCoreTools() {
    // Query tools
    this.registerTool(queryProductsTool);
    this.registerTool(queryPatientsTool);
    this.registerTool(queryAppointmentsTool);

    // Report tools
    this.registerTool(generateQuickReportTool);

    // Action tools
    this.registerTool(scheduleAppointmentTool);
    this.registerTool(createOrderTool);
  }

  registerTool(toolDefinition) {
    const { name, description, inputSchema, handler } = toolDefinition;

    const unregister = navigator.webMCP.registerTool(
      { name, description, inputSchema },
      (params) => this.wrapHandler(handler, params)
    );

    this.registeredTools.set(name, unregister);
    return unregister;
  }

  unregisterTool(name) {
    const unregister = this.registeredTools.get(name);
    if (unregister) {
      unregister();
      this.registeredTools.delete(name);
    }
  }

  wrapHandler(handler, params) {
    // Add context and error handling
    try {
      return handler(params, this.sessionContext);
    } catch (error) {
      console.error(`WebMCP tool error (${handler.name}):`, error);
      return {
        success: false,
        error: error.message || "An unexpected error occurred"
      };
    }
  }

  setupPageLifecycleHandlers() {
    // Enable/disable tools based on page state
    window.addEventListener('page:products', () => {
      this.registerTool(queryProductsTool);
      this.registerTool(createOrderTool);
    });

    window.addEventListener('page:appointments', () => {
      this.registerTool(queryAppointmentsTool);
      this.registerTool(scheduleAppointmentTool);
    });
  }
}

// ============================================
// TOOL DEFINITIONS
// ============================================

// Tool 1: Query Products
const queryProductsTool = {
  name: "queryProducts",
  description: `Search and filter products in the MediSync inventory.
  Returns product details including pricing, stock levels, and categories.
  Use when users want to find specific products or check availability.`,

  inputSchema: {
    type: "object",
    properties: {
      query: {
        type: "string",
        description: "Product name or description search"
      },
      category: {
        type: "string",
        enum: ["medication", "equipment", "supplies", "consumables"],
        description: "Product category filter"
      },
      inStock: {
        type: "boolean",
        description: "Filter to show only in-stock items"
      },
      maxPrice: {
        type: "number",
        description: "Maximum price filter"
      },
      limit: {
        type: "integer",
        minimum: 1,
        maximum: 50,
        default: 20
      }
    }
  },

  async handler(params, context) {
    // Use existing page search functionality
    const productStore = window.__medisync?.productStore;
    if (!productStore) {
      return { success: false, error: "Product data not available" };
    }

    // Build search filters
    const filters = {};
    if (params.category) filters.category = params.category;
    if (params.inStock !== undefined) filters.inStock = params.inStock;
    if (params.maxPrice) filters.maxPrice = params.maxPrice;

    // Execute search
    const results = await productStore.search(params.query || "", {
      filters,
      limit: params.limit || 20
    });

    // Update UI to show results
    productStore.displayResults(results);

    // Return structured data
    return {
      success: true,
      products: results.map(p => ({
        id: p.id,
        name: p.name,
        sku: p.sku,
        category: p.category,
        price: p.price,
        currency: p.currency,
        stockLevel: p.stock,
        inStock: p.stock > 0,
        location: p.warehouseLocation
      })),
      totalCount: results.length,
      query: params.query
    };
  }
};

// Tool 2: Query Patients
const queryPatientsTool = {
  name: "queryPatients",
  description: `Search patient records in the HIMS system.
  Returns patient information based on search criteria.
  Note: Respects privacy settings and user permissions.`,

  inputSchema: {
    type: "object",
    properties: {
      query: {
        type: "string",
        description: "Patient name or ID search"
      },
      status: {
        type: "string",
        enum: ["active", "inactive", "all"],
        default: "active"
      },
      limit: {
        type: "integer",
        minimum: 1,
        maximum: 50,
        default: 20
      }
    }
  },

  async handler(params, context) {
    // Check permissions
    if (!context?.permissions?.canViewPatients) {
      return {
        success: false,
        error: "You don't have permission to view patient records"
      };
    }

    const patientStore = window.__medisync?.patientStore;
    if (!patientStore) {
      return { success: false, error: "Patient data not available" };
    }

    const results = await patientStore.search(params.query || "", {
      status: params.status || "active",
      limit: params.limit || 20
    });

    // Update UI
    patientStore.displayResults(results);

    // Return data (with PII masking based on role)
    return {
      success: true,
      patients: results.map(p => ({
        id: p.id,
        name: context.permissions.canViewPHI ? p.name : maskName(p.name),
        lastVisit: p.lastVisitDate,
        status: p.status
      })),
      totalCount: results.length
    };
  }
};

// Tool 3: Query Appointments
const queryAppointmentsTool = {
  name: "queryAppointments",
  description: `Search and filter appointment schedules.
  Returns upcoming and past appointments based on criteria.`,

  inputSchema: {
    type: "object",
    properties: {
      dateFrom: {
        type: "string",
        format: "date",
        description: "Start date filter"
      },
      dateTo: {
        type: "string",
        format: "date",
        description: "End date filter"
      },
      provider: {
        type: "string",
        description: "Provider/doctor name filter"
      },
      status: {
        type: "string",
        enum: ["scheduled", "completed", "cancelled", "all"],
        default: "scheduled"
      }
    }
  },

  async handler(params, context) {
    const appointmentStore = window.__medisync?.appointmentStore;
    if (!appointmentStore) {
      return { success: false, error: "Appointment data not available" };
    }

    const results = await appointmentStore.query({
      dateFrom: params.dateFrom,
      dateTo: params.dateTo,
      provider: params.provider,
      status: params.status || "scheduled"
    });

    // Update calendar UI
    appointmentStore.highlightAppointments(results);

    return {
      success: true,
      appointments: results.map(a => ({
        id: a.id,
        date: a.date,
        time: a.time,
        duration: a.duration,
        patientId: a.patientId,
        provider: a.providerName,
        type: a.appointmentType,
        status: a.status
      })),
      totalCount: results.length
    };
  }
};

// Tool 4: Generate Quick Report
const generateQuickReportTool = {
  name: "generateQuickReport",
  description: `Generate a quick business report for common metrics.
  Creates visual summaries of key performance indicators.`,

  inputSchema: {
    type: "object",
    properties: {
      reportType: {
        type: "string",
        enum: ["revenue", "appointments", "inventory", "patients"],
        description: "Type of report to generate"
      },
      period: {
        type: "string",
        enum: ["today", "week", "month", "quarter", "year"],
        default: "month"
      }
    },
    required: ["reportType"]
  },

  async handler(params, context) {
    const reportService = window.__medisync?.reportService;
    if (!reportService) {
      return { success: false, error: "Report service not available" };
    }

    const report = await reportService.generateQuickReport(
      params.reportType,
      params.period || "month"
    );

    // Display report in UI
    reportService.showReportModal(report);

    return {
      success: true,
      reportType: params.reportType,
      period: params.period || "month",
      summary: report.summary,
      dataPoints: report.dataPoints,
      generatedAt: new Date().toISOString()
    };
  }
};

// Tool 5: Schedule Appointment (with HITL)
const scheduleAppointmentTool = {
  name: "scheduleAppointment",
  description: `Schedule a new appointment for a patient.
  Shows a form to collect required information.`,

  inputSchema: {
    type: "object",
    properties: {
      patientId: {
        type: "string",
        description: "Patient ID (optional - will prompt if not provided)"
      },
      appointmentType: {
        type: "string",
        description: "Type of appointment"
      },
      preferredDate: {
        type: "string",
        format: "date"
      }
    }
  },

  async handler(params, context) {
    const appointmentService = window.__medisync?.appointmentService;
    if (!appointmentService) {
      return { success: false, error: "Appointment service not available" };
    }

    // If patient not specified, show patient selector
    if (!params.patientId) {
      const selectedPatient = await appointmentService.showPatientSelector();
      if (!selectedPatient) {
        return { success: false, reason: "Patient selection cancelled" };
      }
      params.patientId = selectedPatient.id;
    }

    // Show appointment form to collect details
    const formData = await appointmentService.showAppointmentForm({
      patientId: params.patientId,
      appointmentType: params.appointmentType,
      preferredDate: params.preferredDate
    });

    if (!formData) {
      return { success: false, reason: "Form cancelled by user" };
    }

    // Create the appointment
    const appointment = await appointmentService.create(formData);

    // Update calendar
    appointmentService.refreshCalendar();

    return {
      success: true,
      appointment: {
        id: appointment.id,
        date: appointment.date,
        time: appointment.time,
        patientId: appointment.patientId,
        provider: appointment.providerName
      },
      message: `Appointment scheduled for ${appointment.date} at ${appointment.time}`
    };
  }
};

// Tool 6: Create Order (with confirmation)
const createOrderTool = {
  name: "createOrder",
  description: `Create a new order for products/inventory.
  Requires confirmation before submission.`,

  inputSchema: {
    type: "object",
    properties: {
      items: {
        type: "array",
        items: {
          type: "object",
          properties: {
            productId: { type: "string" },
            quantity: { type: "integer", minimum: 1 }
          },
          required: ["productId", "quantity"]
        },
        description: "List of products to order"
      },
      notes: {
        type: "string",
        description: "Order notes or special instructions"
      }
    },
    required: ["items"]
  },

  async handler(params, context) {
    const orderService = window.__medisync?.orderService;
    if (!orderService) {
      return { success: false, error: "Order service not available" };
    }

    // Fetch product details for confirmation
    const itemDetails = await Promise.all(
      params.items.map(async (item) => {
        const product = await orderService.getProduct(item.productId);
        return {
          ...item,
          name: product.name,
          price: product.price,
          total: product.price * item.quantity
        };
      })
    );

    const orderTotal = itemDetails.reduce((sum, item) => sum + item.total, 0);

    // Show confirmation dialog
    const confirmed = await orderService.showConfirmationDialog({
      title: "Confirm Order",
      items: itemDetails,
      total: orderTotal,
      notes: params.notes
    });

    if (!confirmed) {
      return { success: false, reason: "Order cancelled by user" };
    }

    // Create the order
    const order = await orderService.create({
      items: params.items,
      notes: params.notes,
      total: orderTotal
    });

    // Update UI
    orderService.showOrderConfirmation(order);

    return {
      success: true,
      order: {
        id: order.id,
        itemCount: order.items.length,
        total: order.total,
        status: order.status,
        createdAt: order.createdAt
      }
    };
  }
};

// ============================================
// INITIALIZATION
// ============================================

// Initialize when DOM is ready
document.addEventListener('DOMContentLoaded', () => {
  const webMCPManager = new MediSyncWebMCPManager();

  // Get session context from app
  const context = {
    userId: window.__medisync?.userId,
    companyId: window.__medisync?.companyId,
    permissions: window.__medisync?.permissions || {},
    locale: window.__medisync?.locale || 'en'
  };

  webMCPManager.init(context)
    .then(success => {
      if (success) {
        console.log('WebMCP tools registered successfully');
      }
    })
    .catch(error => {
      console.error('Failed to initialize WebMCP:', error);
    });

  // Export for debugging
  window.__webMCP = webMCPManager;
});

// ============================================
// HELPER FUNCTIONS
// ============================================

function maskName(name) {
  if (!name) return '';
  const parts = name.split(' ');
  return parts.map(part => part[0] + '***').join(' ');
}
