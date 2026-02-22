/**
 * DocumentUploader Component
 *
 * Drag-and-drop file upload component with progress indication.
 * Supports multiple files, file type/size validation, and RTL layout.
 */

import { useState, useCallback, useRef } from 'react'
import { useTranslation } from 'react-i18next'
import { documentApi } from '../../services/documents'
import { LiquidGlassCard } from '../ui/LiquidGlassCard'
import { LiquidGlassButton } from '../ui/LiquidGlassButton'
import type { BulkUploadResponse } from '../../types/documents'

interface DocumentUploaderProps {
  onUploadComplete?: (response: BulkUploadResponse) => void
  onUploadError?: (error: Error) => void
  maxFiles?: number
  maxSizeMB?: number
  className?: string
}

interface FileWithProgress {
  file: File
  progress: number
  status: 'pending' | 'uploading' | 'success' | 'error'
  error?: string
  documentId?: string
}

const ALLOWED_TYPES = [
  'application/pdf',
  'image/jpeg',
  'image/png',
  'image/tiff',
  'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet',
  'text/csv',
]

export function DocumentUploader({
  onUploadComplete,
  onUploadError,
  maxFiles = 50,
  maxSizeMB = 25,
  className = '',
}: DocumentUploaderProps) {
  const { t, i18n } = useTranslation()
  const [files, setFiles] = useState<FileWithProgress[]>([])
  const [isDragging, setIsDragging] = useState(false)
  const [isUploading, setIsUploading] = useState(false)
  const fileInputRef = useRef<HTMLInputElement>(null)

  const isRTL = i18n.language === 'ar'

  const validateFile = useCallback(
    (file: File): string | null => {
      if (!ALLOWED_TYPES.includes(file.type)) {
        return t('documents.upload.invalidType', 'Invalid file type')
      }
      if (file.size > maxSizeMB * 1024 * 1024) {
        return t('documents.upload.tooLarge', 'File exceeds {{max}}MB limit', { max: maxSizeMB })
      }
      return null
    },
    [t, maxSizeMB]
  )

  const handleFiles = useCallback(
    (newFiles: FileList | File[]) => {
      const fileArray = Array.from(newFiles)

      if (files.length + fileArray.length > maxFiles) {
        onUploadError?.(new Error(t('documents.upload.tooMany', 'Maximum {{max}} files allowed', { max: maxFiles })))
        return
      }

      const validatedFiles: FileWithProgress[] = fileArray.map((file) => {
        const error = validateFile(file)
        return {
          file,
          progress: 0,
          status: error ? 'error' : 'pending',
          error: error || undefined,
        }
      })

      setFiles((prev) => [...prev, ...validatedFiles])
    },
    [files.length, maxFiles, validateFile, t, onUploadError]
  )

  const handleDragOver = useCallback((e: React.DragEvent) => {
    e.preventDefault()
    setIsDragging(true)
  }, [])

  const handleDragLeave = useCallback((e: React.DragEvent) => {
    e.preventDefault()
    setIsDragging(false)
  }, [])

  const handleDrop = useCallback(
    (e: React.DragEvent) => {
      e.preventDefault()
      setIsDragging(false)
      handleFiles(e.dataTransfer.files)
    },
    [handleFiles]
  )

  const handleInputChange = useCallback(
    (e: React.ChangeEvent<HTMLInputElement>) => {
      if (e.target.files) {
        handleFiles(e.target.files)
      }
    },
    [handleFiles]
  )

  const removeFile = useCallback((index: number) => {
    setFiles((prev) => prev.filter((_, i) => i !== index))
  }, [])

  const uploadFiles = useCallback(async () => {
    const pendingFiles = files.filter((f) => f.status === 'pending')
    if (pendingFiles.length === 0) return

    setIsUploading(true)

    try {
      // Upload files one by one for progress tracking
      const fileObjects = pendingFiles.map((f) => f.file)
      const response = await documentApi.bulkUpload(fileObjects as unknown as File[])

      // Update file statuses
      setFiles((prev) => {
        const updated = [...prev]
        let pendingIndex = 0
        for (let i = 0; i < updated.length; i++) {
          if (updated[i].status === 'pending') {
            const successId = response.uploadedFiles.find(
              (id) => updated[i].file.name === pendingFiles[pendingIndex]?.file.name
            )
            const failure = response.failedFiles.find(
              (f) => f.filename === updated[i].file.name
            )

            if (successId) {
              updated[i] = {
                ...updated[i],
                status: 'success',
                progress: 100,
                documentId: successId,
              }
            } else if (failure) {
              updated[i] = {
                ...updated[i],
                status: 'error',
                error: failure.error,
              }
            }
            pendingIndex++
          }
        }
        return updated
      })

      onUploadComplete?.(response)
    } catch (error) {
      onUploadError?.(error as Error)
    } finally {
      setIsUploading(false)
    }
  }, [files, onUploadComplete, onUploadError])

  const formatFileSize = (bytes: number): string => {
    if (bytes < 1024) return `${bytes} B`
    if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`
    return `${(bytes / (1024 * 1024)).toFixed(1)} MB`
  }

  return (
    <div className={`document-uploader ${isRTL ? 'rtl' : 'ltr'} ${className}`} dir={isRTL ? 'rtl' : 'ltr'}>
      {/* Drop Zone */}
      <LiquidGlassCard
        intensity="medium"
        hover="shimmer"
        className={`border-2 border-dashed p-8 text-center cursor-pointer ${
          isDragging ? 'border-blue-400/50' : 'border-white/20'
        }`}
        onDragOver={handleDragOver}
        onDragLeave={handleDragLeave}
        onDrop={handleDrop}
        onClick={() => fileInputRef.current?.click()}
      >
        <input
          ref={fileInputRef}
          type="file"
          multiple
          accept=".pdf,.jpg,.jpeg,.png,.tiff,.tif,.xlsx,.csv"
          onChange={handleInputChange}
          className="hidden"
          aria-label={t('documents.upload.selectFiles', 'Select files to upload')}
        />

        <div className="flex flex-col items-center gap-2">
          <div className="text-4xl mb-4">
            {isDragging ? '📥' : '📄'}
          </div>

          <p className="liquid-text-primary font-medium text-lg">
            {t('documents.upload.dropzone', 'Drag and drop files here')}
          </p>
          <p className="liquid-text-secondary text-sm mt-1">
            {t('documents.upload.or', 'or')}{' '}
            <span className="text-blue-500 hover:underline">
              {t('documents.upload.browse', 'browse files')}
            </span>
          </p>
          <p className="liquid-text-muted text-xs mt-2">
            {t('documents.upload.formats', 'PDF, JPEG, PNG, TIFF, XLSX, CSV • Max {{max}}MB', { max: maxSizeMB })}
          </p>
        </div>
      </LiquidGlassCard>

      {/* File List */}
      {files.length > 0 && (
        <div className="mt-4 space-y-2">
          {files.map((fileWithProgress, index) => (
            <LiquidGlassCard
              key={`${fileWithProgress.file.name}-${index}`}
              intensity="subtle"
              hover="lift"
              className="p-3"
            >
              <div className="flex items-center gap-3">
                {/* File Icon */}
                <div className="flex-shrink-0">
                  <div className="text-2xl">📄</div>
                </div>

                {/* File Info */}
                <div className="flex-1 min-w-0">
                  <p className="text-sm font-medium liquid-text-primary truncate">
                    {fileWithProgress.file.name}
                  </p>
                  <p className="text-xs liquid-text-muted">
                    {formatFileSize(fileWithProgress.file.size)}
                  </p>

                  {/* Progress Bar */}
                  {fileWithProgress.status === 'uploading' && (
                    <div className="mt-1 h-1.5 bg-white/10 rounded-full overflow-hidden">
                      <div
                        className="h-full bg-gradient-to-r from-blue-500 to-teal-500 transition-all duration-300"
                        style={{ width: `${fileWithProgress.progress}%` }}
                      />
                    </div>
                  )}

                  {/* Error Message */}
                  {fileWithProgress.status === 'error' && (
                    <p className="text-xs text-red-400 mt-1">{fileWithProgress.error}</p>
                  )}

                  {/* Success Message */}
                  {fileWithProgress.status === 'success' && (
                    <p className="text-xs text-teal-400 mt-1">
                      {t('documents.upload.success', 'Uploaded successfully')}
                    </p>
                  )}
                </div>

                {/* Status Icon */}
                <div className="flex-shrink-0">
                  {fileWithProgress.status === 'success' && (
                    <span className="text-teal-400 text-lg">✓</span>
                  )}
                  {fileWithProgress.status === 'error' && (
                    <span className="text-red-400 text-lg">✕</span>
                  )}
                  {fileWithProgress.status === 'pending' && !isUploading && (
                    <button
                      onClick={(e) => {
                        e.stopPropagation()
                        removeFile(index)
                      }}
                      className="liquid-text-muted hover:text-red-400 transition-colors"
                      aria-label={t('documents.upload.remove', 'Remove file')}
                    >
                      <span className="text-lg">✕</span>
                    </button>
                  )}
                </div>
              </div>
            </LiquidGlassCard>
          ))}
        </div>
      )}

      {/* Upload Button */}
      {files.some((f) => f.status === 'pending') && (
        <div className="mt-4 flex justify-end">
          <LiquidGlassButton
            variant="primary"
            onClick={uploadFiles}
            disabled={isUploading}
            isLoading={isUploading}
          >
            {isUploading
              ? t('documents.upload.uploading', 'Uploading...')
              : t('documents.upload.start', 'Start Upload')}
          </LiquidGlassButton>
        </div>
      )}
    </div>
  )
}

export default DocumentUploader
