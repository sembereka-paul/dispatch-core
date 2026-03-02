import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "~/components/ui/card";
import { useEffect, useLayoutEffect, useState } from "react";
import { Heart, Repeat, Reply } from "lucide-react";

export type Account = {
  id: string;
  username?: string;
  acct?: string;
  display_name?: string;
  url?: string;
  avatar?: string;
  avatar_static?: string;
  header?: string;
  header_static?: string;
  locked?: boolean;
  bot?: boolean;
  created_at?: string;
  note?: string;
};

export type MediaAttachment = {
  id: string;
  type: string;
  url?: string;
  preview_url?: string;
  remote_url?: string | null;
  text_url?: string;
  description?: string | null;
};

export type Mention = {
  id: string;
  username: string;
  url: string;
  acct: string;
};

export type Tag = {
  name: string;
  url: string;
};

export type Application = {
  name?: string;
  website?: string | null;
};
export type Post = {
  id: string;
  created_at: string;
  in_reply_to_id: string | null;
  in_reply_to_account_id: string | null;
  sensitive: boolean;
  spoiler_text: string;
  visibility: 'public' | 'unlisted' | 'private' | 'direct' | string;
  language: string | null;
  uri: string;
  url: string | null;
  replies_count?: number;
  reblogs_count?: number;
  favourites_count?: number;
  favourited?: boolean | null;
  reblogged?: boolean | null;
  muted?: boolean | null;
  bookmarked?: boolean | null;
  pinned?: boolean | null;
  content?: string;
  filtered?: any[];
  account?: Account;
  media_attachments?: MediaAttachment[];
  mentions?: Mention[];
  tags?: Tag[];
  application?: Application | null;
};
export default function Home() {
  const data = useEventSourcing()

  useEffect(() => {
    console.log(data, 'hehe')
  }, [data])

  return <>
    <div>
      <div className="flex flex-cols py-4 justify-center">

        {data &&
          <Card size="sm" className="mx-auto w-full max-w-sm">
            <CardHeader>
              <CardTitle>{data.account?.display_name}</CardTitle>
              <CardDescription>@{data.account?.username}</CardDescription>
              {/* <CardAction>Card Action</CardAction> */}
            </CardHeader>
            <CardContent>
              <div dangerouslySetInnerHTML={{ __html: data.content! }} />
            </CardContent>
            <CardFooter className="flex flex-cols justify-between gap-2">
              <div>
                <span className="flex flex-cols align-center gap-1">
                  <Reply size={22}></Reply>
                  {data.replies_count}
                </span>

              </div>
              <div>
                <span className="flex flex-cols align-center gap-1">
                  <Repeat size={22}></Repeat>
                  {data.reblogs_count}
                </span>
              </div>
              <div>
                <span className="flex flex-cols align-center gap-1">
                  <Heart size={22}></Heart>
                  {data.favourites_count}
                </span>
              </div>
            </CardFooter>
          </Card>
        }
      </div>
    </div>
  </>
}


const useEventSourcing = () => {
  const [data, setData] = useState<Post | undefined>(undefined)
  useLayoutEffect(() => {
    const eventSource = new EventSource('http://localhost:8080/notifications/new');

    // Listen for the default "message" event
    eventSource.onmessage = (event) => {
      console.log('Received:', event.data);
      const data = JSON.parse(event.data) as Post
      setData(data)
      // Process the data
    };

    // Handle connection opened
    eventSource.onopen = (event) => {
      console.log('Connection established');
    };

    // Handle errors (including disconnections)
    eventSource.onerror = (event) => {
      console.error('EventSource error:', event);
      // Browser will automatically attempt to reconnect
    };

  }, [])

  return data
  // Create a connection to the SSE endpoint
}
