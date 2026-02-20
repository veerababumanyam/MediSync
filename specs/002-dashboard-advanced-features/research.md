# Research: Dashboard, Chat UI & Advanced Features with i18n

**Feature**: 002-dashboard-advanced-features
**Date**: 2026-02-20
**Status**: Phase 0 Complete

---

## 1. CopilotKit Integration for Streaming AI Responses

### Decision
Use CopilotKit's `useCopilotChat` hook with `CopilotPopup` or custom `CopilotChat` component for streaming natural language responses with inline chart rendering.

### Rationale
- **Native streaming support**: CopilotKit handles SSE/WebSocket streaming out of the box
- **Generative UI**: Can render React components dynamically via `renderDataStream` callback
- **React 19 compatible**: Version 1.3.6 supports React 19's concurrent features
- **MIT licensed**: Meets open-source requirement

### Implementation Pattern
```typescript
// Chat interface with streaming
import { useCopilotChat, CopilotMessage } from "@copilotkit/react-core";

function ChatInterface() {
  const { messages, appendMessage, isLoading } = useCopilotChat({
    id: "medisync-chat",
    instructions: "You are a healthcare BI assistant...",
  });

  // Streaming response renders progressively
  return (
    <div>
      {messages.map((msg) => (
        <CopilotMessage key={msg.id} message={msg}>
          {msg.content.type === "chart" && (
            <ChartRenderer spec={msg.content.chartSpec} />
          )}
        </CopilotMessage>
      ))}
    </div>
  );
}
```

### Alternatives Considered
| Alternative | Rejected Because |
|-------------|------------------|
| **Vercel AI SDK** | Uses React Server Components, requires Next.js; our stack is Vite |
| **Custom SSE handler** | More boilerplate, no built-in UI primitives |
| **LangChain.js** | Heavier, more suited for complex agent chains than chat UI |

### References
- [CopilotKit Documentation](https://docs.copilotkit.ai/)
- [React 19 + CopilotKit Guide](https://docs.copilotkit.ai/reference/react/19)

---

## 2. Apache ECharts RTL and Arabic Locale Support

### Decision
Use Apache ECharts with built-in Arabic locale (`echarts/lang/ar`) and CSS `direction: rtl` for chart container mirroring.

### Rationale
- **Native Arabic support**: ECharts includes `lang/ar.js` with Arabic UI strings
- **RTL handling**: Works with CSS logical properties; no manual axis reversal needed
- **Apache 2.0 license**: Meets open-source requirement
- **Already in stack**: Project uses `echarts` + `echarts-for-react`

### Implementation Pattern
```typescript
import * as echarts from 'echarts';
import 'echarts/lang/ar'; // Arabic locale

function ChartRenderer({ spec, locale }: { spec: ChartSpec; locale: string }) {
  const chartRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const chart = echarts.init(chartRef.current, null, {
      locale: locale === 'ar' ? 'AR' : 'EN',
    });
    chart.setOption(spec);
    return () => chart.dispose();
  }, [spec, locale]);

  return (
    <div
      ref={chartRef}
      dir={locale === 'ar' ? 'rtl' : 'ltr'}
      className="w-full h-80"
    />
  );
}
```

### RTL Considerations
- **Legend position**: Use `legend.left: 'right'` for Arabic
- **Axis labels**: ECharts auto-handles with locale
- **Tooltips**: Ensure tooltip container has `dir` attribute

### Alternatives Considered
| Alternative | Rejected Because |
|-------------|------------------|
| **Recharts** | Less mature RTL support, smaller community |
| **Chart.js** | Requires manual RTL workarounds |
| **D3.js** | Too low-level, would need custom RTL handling |

---

## 3. PDF Generation with Arabic/RTL Support

### Decision
Use **Puppeteer** (via `puppeteer-core`) with a headless Chrome instance to generate PDFs from HTML, ensuring proper Arabic font rendering and RTL layout.

### Rationale
- **Native RTL support**: Chrome renders Arabic perfectly with `dir="rtl"`
- **Font control**: Can load system Arabic fonts (Amiri, Noto Arabic)
- **WYSIWYG**: PDF matches on-screen rendering exactly
- **MIT licensed** (puppeteer-core)

### Implementation Pattern
```go
// Backend PDF generation
func (s *ReportService) GeneratePDF(ctx context.Context, report Report) ([]byte, error) {
    // 1. Render HTML template with Arabic content
    html, err := s.renderTemplate(report, locale)
    if err != nil {
        return nil, err
    }

    // 2. Use browserless.io or local Chrome
    browser, err := puppeteer.Launch(puppeteer.LaunchOptions{
        Headless: true,
        Args: []string{"--font-render-hinting=none"},
    })
    defer browser.Close()

    page, _ := browser.NewPage()
    page.SetContent(html)
    page.Evaluate(fmt.Sprintf(`document.dir = "%s"`, locale.Direction))

    return page.PDF(puppeteer.PDFOptions{
        Format: "A4",
        PrintBackground: true,
    })
}
```

### Font Requirements
- Install `fonts-noto-core` or `fonts-arabeyes` on server
- CSS: `font-family: 'Amiri', 'Noto Sans Arabic', sans-serif;`

### Alternatives Considered
| Alternative | Rejected Because |
|-------------|------------------|
| **gofpdf** | Limited Unicode/Arabic support, requires manual glyph shaping |
| **wkhtmltopdf** | Older WebKit, inconsistent Arabic rendering |
| **ReportLab (Python)** | Requires Python dependency, licensing for commercial use |

---

## 4. Spreadsheet Export with RTL Support

### Decision
Use **ExcelJS** (Node.js) for `.xlsx` generation with RTL worksheet direction, or pure Go with `excelize` for server-side generation.

### Rationale
- **RTL sheet direction**: ExcelJS supports `views: [{ rightToLeft: true }]`
- **No external dependencies**: Run in Go service
- **Apache 2.0 licensed**

### Implementation Pattern
```go
import "github.com/xuri/excelize/v2"

func (s *ExportService) GenerateSpreadsheet(data TableData, locale string) ([]byte, error) {
    f := excelize.NewFile()
    sheet := "Sheet1"

    // Set RTL for Arabic locale
    if locale == "ar" {
        f.SetSheetView(sheet, &excelize.ViewOptions{
            RightToLeft: boolPtr(true),
        })
    }

    // Write data
    for i, row := range data.Rows {
        for j, cell := range row.Cells {
            coord, _ := excelize.CoordinatesToCellName(j+1, i+1)
            f.SetCellValue(sheet, coord, cell.Value)
        }
    }

    return f.WriteToBuffer()
}
```

### Alternatives Considered
| Alternative | Rejected Because |
|-------------|------------------|
| **CSV** | No RTL support; use only for data interchange |
| **golang/spreadsheet** | Less mature, limited RTL support |
| **LibreOffice UNO** | Requires full LibreOffice installation |

---

## 5. WebSocket Streaming in Go Backend

### Decision
Use **gorilla/websocket** (already in dependencies) with `chi` middleware for WebSocket upgrade. Implement message streaming from AI agent responses.

### Rationale
- **Already in stack**: `github.com/gorilla/websocket v1.5.3`
- **Chi integration**: Works with existing router
- **BSD-3 licensed**

### Implementation Pattern
```go
// internal/api/websocket/stream.go
var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        return true // Configure for production
    },
}

func (h *StreamHandler) HandleChatStream(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        return
    }
    defer conn.Close()

    // Get user locale from context (set by middleware)
    locale := ctx.Value("locale").(string)

    // Stream from AI agent
    stream, err := h.agentService.TextToSQLStream(r.Context(), query, locale)
    if err != nil {
        conn.WriteJSON(StreamError{Error: err.Error()})
        return
    }

    for chunk := range stream {
        conn.WriteJSON(StreamChunk{
            Type: chunk.Type,    // "text", "chart", "table"
            Content: chunk.Data,
            Done: chunk.Done,
        })
    }
}
```

### Message Types
```typescript
interface StreamChunk {
  type: 'text' | 'chart' | 'table' | 'error';
  content: string | ChartSpec | TableData;
  done: boolean;
}
```

---

## 6. Alert Scheduling and Evaluation

### Decision
Use **NATS JetStream** (already in stack) with delayed messages for alert scheduling, plus a dedicated alert evaluator service.

### Rationale
- **Already in stack**: `github.com/nats-io/nats.go v1.39.1`
- **Delayed delivery**: JetStream supports `DeliverDelay`
- **Durable queues**: Survives restarts
- **Apache 2.0 licensed**

### Implementation Pattern
```go
// Alert scheduling on rule creation
func (s *AlertService) ScheduleAlert(ctx context.Context, rule AlertRule) error {
    js, _ := s.nc.JetStream()

    // Schedule periodic evaluation
    _, err := js.AddConsumer("ALERTS", &nats.ConsumerConfig{
        Durable:   fmt.Sprintf("alert-%d", rule.ID),
        DeliverSubject: fmt.Sprintf("alert.evaluate.%d", rule.ID),
        DeliverDelay:   rule.CheckInterval, // e.g., 5 * time.Minute
    })
    return err
}

// Alert evaluator
func (s *AlertService) EvaluateAlert(ctx context.Context, rule AlertRule) error {
    // 1. Execute metric query
    value, err := s.metricStore.QueryValue(ctx, rule.MetricID)
    if err != nil {
        return err
    }

    // 2. Check threshold
    triggered := s.compareThreshold(value, rule.Operator, rule.Threshold)

    if triggered {
        // 3. Send notification
        s.notifier.Send(ctx, Notification{
            UserID:  rule.UserID,
            Message: s.formatMessage(rule, value, locale),
            Channels: rule.Channels,
        })
    }
    return nil
}
```

### Alternatives Considered
| Alternative | Rejected Because |
|-------------|------------------|
| **Cron** | Less flexible for dynamic scheduling |
| **Temporal** | Overkill for simple alert evaluation |
| **Database polling** | Less efficient than message-driven |

---

## 7. Report Scheduling and Delivery

### Decision
Use **NATS JetStream** for scheduling + **SMTP** (via `net/smtp`) for email delivery. Store report configuration and history in PostgreSQL.

### Rationale
- **Native scheduling**: JetStream delayed messages
- **Retry logic**: Built-in redelivery with backoff
- **Audit trail**: All deliveries logged to `scheduled_report_runs` table

### Implementation Pattern
```go
// Schedule recurring report
func (s *ReportService) ScheduleReport(ctx context.Context, report ScheduledReport) error {
    // Calculate next run time
    nextRun := s.calculateNextRun(report.Schedule)

    // Schedule via NATS
    js, _ := s.nc.JetStream()
    _, err := js.PublishMsg(&nats.Msg{
        Subject: "report.generate",
        Data:    marshalJSON(report),
        Header: nats.Header{
            "Nats-Delay": []string{time.Until(nextRun).String()},
        },
    })
    return err
}

// Generate and send report
func (s *ReportService) GenerateAndSend(ctx context.Context, report ScheduledReport) error {
    // 1. Execute report query
    data, err := s.queryExecutor.Execute(ctx, report.QueryID)

    // 2. Generate file
    var attachment []byte
    switch report.Format {
    case "pdf":
        attachment, err = s.pdfGenerator.Generate(ctx, data, report.Locale)
    case "xlsx":
        attachment, err = s.spreadsheetGenerator.Generate(ctx, data, report.Locale)
    case "csv":
        attachment, err = s.csvGenerator.Generate(ctx, data)
    }

    // 3. Send email
    return s.emailSender.Send(ctx, Email{
        To:        report.Recipients,
        Subject:   s.t("reports.subject", report.Locale),
        Body:      s.t("reports.body", report.Locale),
        Attachment: attachment,
    })
}
```

---

## 8. Locale Detection and Routing

### Decision
Follow priority chain: `user_preferences.locale` (JWT claim) → `Accept-Language` header → URL param `?lang=` → default `en`.

### Rationale
- **Consistent with CLAUDE.md**: Matches documented i18n pattern
- **JWT-first**: User preference stored in Keycloak, included in token
- **Fallback chain**: Ensures locale always determined

### Implementation Pattern
```go
// internal/api/middleware/locale.go
func LocaleMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        var locale string

        // 1. Check JWT claim (set by auth middleware)
        if jwtLocale := r.Context().Value("user_locale"); jwtLocale != nil {
            locale = jwtLocale.(string)
        }

        // 2. Check Accept-Language header
        if locale == "" {
            acceptLang := r.Header.Get("Accept-Language")
            if strings.HasPrefix(acceptLang, "ar") {
                locale = "ar"
            }
        }

        // 3. Check URL param
        if locale == "" {
            locale = r.URL.Query().Get("lang")
        }

        // 4. Default
        if locale == "" {
            locale = "en"
        }

        // Set in context
        ctx := context.WithValue(r.Context(), "locale", locale)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
```

---

## 9. Frontend State Management

### Decision
Use **React Query** (`@tanstack/react-query`) for server state (API data) and **Zustand** for client state (UI preferences, sidebar state).

### Rationale
- **Separation of concerns**: Server state vs client state
- **Caching**: React Query handles cache, invalidation, background refetch
- **Lightweight**: Zustand is minimal boilerplate
- **Both MIT licensed**

### Implementation Pattern
```typescript
// Server state with React Query
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';

export function usePinnedCharts() {
  return useQuery({
    queryKey: ['pinned-charts'],
    queryFn: () => api.get('/api/dashboard/charts'),
    staleTime: 5 * 60 * 1000, // 5 minutes
  });
}

export function usePinChart() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (chart: PinChartRequest) => api.post('/api/dashboard/charts', chart),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['pinned-charts'] });
    },
  });
}

// Client state with Zustand
import { create } from 'zustand';

interface UIState {
  sidebarOpen: boolean;
  locale: 'en' | 'ar';
  toggleSidebar: () => void;
  setLocale: (locale: 'en' | 'ar') => void;
}

export const useUIStore = create<UIState>((set) => ({
  sidebarOpen: true,
  locale: 'en',
  toggleSidebar: () => set((s) => ({ sidebarOpen: !s.sidebarOpen })),
  setLocale: (locale) => set({ locale }),
}));
```

### Alternatives Considered
| Alternative | Rejected Because |
|-------------|------------------|
| **Redux Toolkit** | More boilerplate; overkill for this scope |
| **SWR** | Good, but React Query has better mutation support |
| **Jotai/Recoil** | Less ecosystem support for server state |

---

## 10. Mobile Offline Dashboard (Out of Scope for This Phase)

### Research Note
Mobile offline support will use **PowerSync** for local SQLite cache with sync-on-reconnect. Flutter app will read from same `pinned_charts` API endpoint. This is documented for future reference but not implemented in this feature.

### Key Considerations for Future
- PowerSync schema must match PostgreSQL `pinned_charts` table
- Last-updated timestamps for conflict resolution
- Background sync on app resume

---

## Summary of Decisions

| Area | Decision | License |
|------|----------|---------|
| AI Chat UI | CopilotKit | MIT |
| Charts | Apache ECharts with Arabic locale | Apache-2.0 |
| PDF Generation | Puppeteer (headless Chrome) | MIT |
| Spreadsheet Export | excelize (Go) | Apache-2.0 |
| WebSocket | gorilla/websocket | BSD-3 |
| Alert Scheduling | NATS JetStream delayed messages | Apache-2.0 |
| Report Scheduling | NATS JetStream + net/smtp | Apache-2.0 |
| Frontend State | React Query + Zustand | MIT |

All decisions align with:
- ✅ Constitution Principle IV (Open Source Only)
- ✅ Constitution Principle III (i18n by Default)
- ✅ Existing technology stack in `CLAUDE.md`

---

*Research completed: 2026-02-20*
