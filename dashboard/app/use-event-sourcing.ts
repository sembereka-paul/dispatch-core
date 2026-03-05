import { useEffect, useLayoutEffect, useState } from "react"
import type { Post } from "./schema"

export const useEventSourcing = (tag?: string) => {
  const [sub, setSub] = useState<string | undefined>(tag)
  const [data, setData] = useState<Post | undefined>(undefined)
  const [error, setError] = useState<Error | undefined>(undefined)

  const trackTag = (name?: string) => setSub(name)

  useEffect(() => { setSub(tag) }, [tag])

  useLayoutEffect(() => {
    if (!sub) return

    const eventSource = new EventSource(`http://localhost:8080/notifications/${sub}`);

    // Listen for the default "message" event
    eventSource.onmessage = (event) => {
      console.log('Received:', event.data);
      const data = JSON.parse(event.data) as Post
      setData(data)
    };

    // Handle connection opened
    eventSource.onopen = (event) => {
      console.log('Connection established');

    };

    // Handle errors (including disconnections)
    eventSource.onerror = (event) => {
      console.error('EventSource error:', event);
      setError(error)
    };

    return () => {
      eventSource.close()
    }

  }, [sub])

  return {
    data,
    trackTag
  }
}
