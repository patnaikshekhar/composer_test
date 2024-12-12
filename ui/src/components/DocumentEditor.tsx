import { useState } from 'react'
import { Button } from "@/components/ui/button"
import { Textarea } from "@/components/ui/textarea"
import { Switch } from "@/components/ui/switch"
import { Label } from "@/components/ui/label"
import ReactMarkdown from 'react-markdown'

// Mock data for initial document content
const initialDocument = `# Welcome to Composer

This is a sample document. You can edit this content and see the changes reflected in real-time.

## Features

- AI-assisted writing
- Real-time markdown preview
- Collaborative editing

Start typing to create your document!
`

export default function DocumentEditor() {
  const [document, setDocument] = useState(initialDocument)
  const [isEditing, setIsEditing] = useState(false)

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
            <Textarea
              value={document}
              onChange={(e) => setDocument(e.target.value)}
              className="w-full h-[calc(100vh-200px)] p-4 border rounded font-mono"
            />
          ) : (
            <div className="prose dark:prose-invert max-w-none">
              <ReactMarkdown>{document}</ReactMarkdown>
            </div>
          )}
        </div>
      </div>
    </div>
  )
}
