import React from "react";
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "~/components/ui/card";
import { useEffect, useLayoutEffect, useState } from "react";
import { Button } from "~/components/ui/button"
import { Field } from "~/components/ui/field"
import { Input } from "~/components/ui/input"
import { BrushCleaning, Heart, Repeat, Reply } from "lucide-react";
import { useEventSourcing } from "~/use-event-sourcing";
import type { Post } from "~/schema";

export default function Home() {
  const [tag, setTag] = useState<string | undefined>()
  const { data, trackTag } = useEventSourcing()
  const [posts, setPosts] = useState<Post[]>([])

  const [tracked, setTracked] = useState(new Set<string>())

  useEffect(() => {
    if (!data) return
    setPosts(old => [data, ...old])
  }, [data])

  useEffect(() => {
    if (!tracked.size) setPosts([])
  }, [tracked])

  const handleTracking = () => {
    if (!tag) return
    trackTag(tag)
    setTracked(new Set([...tracked, tag]))
  }

  const handleClear = () => {
    setTag('')
    setTracked(new Set())
  }

  const renderNowTracking = () => {
    if (!tracked.size) return

    let tags = ''
    tracked.forEach(v => tags = tags + ` #${v}`)

    return <div className="pt-2 text-center text-gray-500" >
      <span>now tracking {tags?.trim()}</span>
    </div>
  }

  return <>
    <div>
      <div className="flex py-4 mb-4 justify-center">
        <div className="w-80">
          <Field orientation="horizontal">
            <Input onChange={(e) => setTag(e.target.value)} value={tag} type="text" placeholder="Enter tag..." />
            <Button onClick={handleTracking} disabled={!tag} className="cursor-pointer">Track</Button>
            <Button disabled={!tracked.size && !tag} onClick={handleClear} className="cursor-pointer text-red-700 bg-red-100">clear</Button>
          </Field>
          {renderNowTracking()}
        </div>
      </div>
      <div className="py-4 justify-center">
        {posts.map(post => {
          return <React.Fragment key={post.id}>
            < Card size="sm" className="py-2 my-2 mx-auto w-full max-w-sm">
              <CardHeader>
                <CardTitle>{post.account?.display_name}</CardTitle>
                <CardDescription>@{post.account?.username}</CardDescription>
              </CardHeader>
              <CardContent>
                <div dangerouslySetInnerHTML={{ __html: post.content! }} />
              </CardContent>
              <CardFooter className="flex flex-cols justify-between gap-2">
                <div>
                  <span className="flex flex-cols align-center gap-1">
                    <Reply size={22}></Reply>
                    {post.replies_count}
                  </span>

                </div>
                <div>
                  <span className="flex flex-cols align-center gap-1">
                    <Repeat size={22}></Repeat>
                    {post.reblogs_count}
                  </span>
                </div>
                <div>
                  <span className="flex flex-cols align-center gap-1">
                    <Heart size={22}></Heart>
                    {post.favourites_count}
                  </span>
                </div>
              </CardFooter>
            </Card>
          </React.Fragment>
        })
        }
        {posts.length === 0 && <div className="w-80 m-10 m-auto">
          <div className="flex justify-center">
            <BrushCleaning></BrushCleaning>

          </div>
          <div className="text-center">
            Empty posts
          </div>
        </div>
        }
      </div>
    </div >
  </>
}
