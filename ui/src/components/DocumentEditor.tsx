'use client'

import { useState } from 'react'
import { Button } from "@/components/ui/button"
import { Switch } from "@/components/ui/switch"
import { Label } from "@/components/ui/label"
import ReactMarkdown from 'react-markdown'
import { Prism as SyntaxHighlighter } from 'react-syntax-highlighter'
import { atomDark } from 'react-syntax-highlighter/dist/esm/styles/prism'
import Editor from '@monaco-editor/react'

interface DocumentEditorProps {
  document: string;
  onDocumentChange: (newDocument: string) => void;
}

export default function DocumentEditor({ document, onDocumentChange }: DocumentEditorProps) {
  const [isEditing, setIsEditing] = useState(false)
  const handleEditorChange = (value: string | undefined) => {
    if (value !== undefined) {
      onDocumentChange(value)
    }
  }

  return (
    <div className="flex-grow p-8 overflow-auto">
      <div className="max-w-3xl mx-auto">
        <div className="mb-4 flex items-center justify-between">
          <div className="flex items-center space-x-2">
            <Switch
              id="edit-mode"
              checked={isEditing}
              onCheckedChange={setIsEditing}
            />
            <Label htmlFor="edit-mode">{isEditing ? 'Editing' : 'Previewing'}</Label>
          </div>
          <Button variant="outline">Save</Button>
        </div>
        <div className="relative">
          {isEditing ? (
            <Editor
              height="calc(100vh - 200px)"
              defaultLanguage="markdown"
              value={document}
              theme="vs-dark"
              onChange={handleEditorChange}
              options={{
                minimap: { enabled: false },
                fontSize: 14,
                lineNumbers: 'on',
                wordWrap: 'on',
                wrappingIndent: 'indent',
                lineDecorationsWidth: 10,
                lineNumbersMinChars: 0,
                glyphMargin: false,
                folding: false,
                scrollBeyondLastLine: false,
                automaticLayout: true,
              }}
            />
          ) : (
            <div className="prose dark:prose-invert max-w-none">
              <ReactMarkdown
                components={{
                  code({ node, inline, className, children, ...props }) {
                    const match = /language-(\w+)/.exec(className || '')
                    return !inline && match ? (
                      <SyntaxHighlighter
                        language={match[1]}
                        style={atomDark}
                        PreTag="div"
                        {...props}
                      >
                        {String(children).replace(/\n$/, '')}
                      </SyntaxHighlighter>
                    ) : (
                      <code className={className} {...props}>
                        {children}
                      </code>
                    )
                  }
                }}
              >
                {document}
              </ReactMarkdown>
            </div>
          )}
        </div>
      </div>
    </div>
  )
}

