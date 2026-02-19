---
name: ocr-extraction
description: Extract structured financial data from multi-modal documents (PDF, Images, Scans). Use for digitizing invoices, bills, and bank statements with high precision and confidence scoring.
---

# OCR Extraction Skill

Guidelines for extracting financial entities from various document formats using the MediSync AI Accountant OCR pipeline.

## Extraction Strategy by Format

### Digital PDFs
- Use **PyMuPDF** to extract the native text layer.
- **Principle**: Preferred over OCR for speed and perfect accuracy when a text layer is available.
- **Failover**: If the PDF is scanned (no text layer), route to the Scanned/Image pipeline.

### Scanned Documents / Images
- Use **PaddleOCR** for layout analysis and text recognition.
- **Table Extraction**: Use `unstructured.io` to segment line items and headers from complex invoice tables.
- **Preprocessing**: Apply grayscale conversion, noise reduction, and deskewing for low-quality mobile scans.

### Handwritten Scans
- Use **PaddleOCR (HTR model)** combined with LLM post-processing.
- **LLM Prompting**: Provide the noisy OCR text to Gemini and ask it to "Correct spelling and extract structured financial fields based on the surrounding context."

## Tool Chain Patterns

### Unstructured.io Segmentation
```python
from unstructured.partition.pdf import partition_pdf

elements = partition_pdf(
    filename="invoice.pdf",
    strategy="hi_res",  # Uses detectron2 for layout analysis
    model_name="yolox"
)

# Filter for Table elements to extract line items
tables = [el for el in elements if el.category == "Table"]
```

### Pydantic Extraction Model
Define the source of truth for extracted fields:
```python
from pydantic import BaseModel, Field
from typing import List, Optional

class LineItem(BaseModel):
    description: str
    quantity: float
    unit_price: float
    total: float

class InvoiceData(BaseModel):
    vendor_name: str
    invoice_number: str
    date: str
    subtotal: float
    tax_amount: float
    total_amount: float
    currency: str = "INR"
    line_items: List[LineItem]
```

## Accuracy & Verification

### Confidence Scoring
- Assign a confidence score (0-1) to every extracted field.
- **Flagging**: If any "Critical Field" (`total_amount`, `vendor`, `date`) is `< 0.85`, the document must be routed to the HITL gate.

### Logic Checks
- Subtotal + Tax = Total.
- Line Item Sum = Subtotal.
- Use these mathematical invariants to increase extraction confidence.

## Accessibility Checklist
- [ ] Provide "Raw Text" view for users to compare against original image.
- [ ] Highlight extracted fields on the original document preview in the UI.
- [ ] Support large file uploads (up to 50MB) and batch processing.
