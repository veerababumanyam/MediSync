/**
 * DocumentUploader Component
 *
 * Drag-and-drop file upload component with progress indication.
 * Supports multiple files, file type/size validation, and RTL layout.
 */

import { useState, useCallback, useRef } from 'react'
import { useTranslation } from 'react-i18next'
import { documentApi } from '../../services/documents'
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
      <div
        className={`
          border-2 border-dashed rounded-lg p-8 text-center cursor-pointer transition-colors
          ${isDragging ? 'border-blue-500 bg-blue-50' : 'border-gray-300 hover:border-gray-400'}
          dark:border-gray-600 dark:hover:border-gray-500
          ${isDragging ? 'dark:bg-blue-900/20' : ''}
        `}
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
          <svg
            className={`w-12 h-12 ${isDragging ? 'text-blue-500' : 'text-gray-400'}`}
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12"
            />
          </svg>

          <p className="text-lg font-medium text-gray-700 dark:text-gray-200">
            {t('documents.upload.dropzone', 'Drag and drop files here')}
          </p>
          <p className="text-sm text-gray-500 dark:text-gray-400">
            {t('documents.upload.or', 'or')}{' '}
            <span className="text-blue-600 dark:text-blue-400 hover:underline">
              {t('documents.upload.browse', 'browse files')}
            </span>
          </p>
          <p className="text-xs text-gray-400 dark:text-gray-500 mt-2">
            {t('documents.upload.formats', 'PDF, JPEG, PNG, TIFF, XLSX, CSV • Max {{max}}MB', { max: maxSizeMB })}
          </p>
        </div>
      </div>

      {/* File List */}
      {files.length > 0 && (
        <div className="mt-4 space-y-2">
          {files.map((fileWithProgress, index) => (
            <div
              key={`${fileWithProgress.file.name}-${index}`}
              className="flex items-center gap-3 p-3 bg-gray-50 dark:bg-gray-800 rounded-lg"
            >
              {/* File Icon */}
              <div className="flex-shrink-0">
                <svg className="w-8 h-8 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
                  />
                </svg>
              </div>

              {/* File Info */}
              <div className="flex-1 min-w-0">
                <p className="text-sm font-medium text-gray-900 dark:text-gray-100 truncate">
                  {fileWithProgress.file.name}
                </p>
                <p className="text-xs text-gray-500 dark:text-gray-400">
                  {formatFileSize(fileWithProgress.file.size)}
                </p>

                {/* Progress Bar */}
                {fileWithProgress.status === 'uploading' && (
                  <div className="mt-1 h-1.5 bg-gray-200 dark:bg-gray-700 rounded-full overflow-hidden">
                    <div
                      className="h-full bg-blue-500 transition-all duration-300"
                      style={{ width: `${fileWithProgress.progress}%` }}
                    />
                  </div>
                )}

                {/* Error Message */}
                {fileWithProgress.status === 'error' && (
                  <p className="text-xs text-red-500 mt-1">{fileWithProgress.error}</p>
                )}

                {/* Success Message */}
                {fileWithProgress.status === 'success' && (
                  <p className="text-xs text-green-500 mt-1">
                    {t('documents.upload.success', 'Uploaded successfully')}
                  </p>
                )}
              </div>

              {/* Status Icon */}
              <div className="flex-shrink-0">
                {fileWithProgress.status === 'success' && (
                  <svg className="w-5 h-5 text-green-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                  </svg>
                )}
                {fileWithProgress.status === 'error' && (
                  <svg className="w-5 h-5 text-red-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                  </svg>
                )}
                {fileWithProgress.status === 'pending' && !isUploading && (
                  <button
                    onClick={(e) => {
                      e.stopPropagation()
                      removeFile(index)
                    }}
                    className="text-gray-400 hover:text-red-500"
                    aria-label={t('documents.upload.remove', 'Remove file')}
                  >
                    <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                    </svg>
                  </button>
                )}
              </div>
            </div>
          ))}
        </div>
      )}

      {/* Upload Button */}
      {files.some((f) => f.status === 'pending') && (
        <div className="mt-4 flex justify-end">
          <button
            onClick={uploadFiles}
            disabled={isUploading}
            className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
          >
            {isUploading
              ? t('documents.upload.uploading', 'Uploading...')
              : t('documents.upload.start', 'Start Upload')}
          </button>
        </div>
      )}
    </div>
  )
}

export default DocumentUploader
