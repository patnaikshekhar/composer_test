import { useState } from 'react'
import ChatSidebar from './components/ChatSidebar'
import DocumentEditor from './components/DocumentEditor'
import { Button } from "@/components/ui/button"
import { MessageSquare } from 'lucide-react'

export default function App() {
  const [isSidebarOpen, setIsSidebarOpen] = useState(true)

  return (
    <div className="flex h-screen bg-background text-foreground">
      {isSidebarOpen && <ChatSidebar />}
      <div className="flex-1 flex flex-col">
        <header className="flex justify-between items-center p-4 border-b">
          <Button variant="ghost" onClick={() => setIsSidebarOpen(!isSidebarOpen)}>
            <MessageSquare className="h-[1.2rem] w-[1.2rem]" />
          </Button>
          <h1 className="text-2xl font-bold">Composer</h1>
          <div />
        </header>
        <DocumentEditor />
      </div>
    </div>
  )
}