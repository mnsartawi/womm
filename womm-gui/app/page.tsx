"use client"
import { motion } from "framer-motion"

import Image from "next/image"
import { Button } from "@/components/ui/button"
import { useEffect, useState } from "react"
import { SiGithub } from "react-icons/si"
import { CogIcon, FileEdit, FilePlus, Globe, ScrollText } from "lucide-react"
import { open } from "@tauri-apps/plugin-shell"
import { useRouter } from "next/navigation"
import {
  Command,
  CommandDialog,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
} from "@/components/ui/command"

export default function Home() {
  const router = useRouter()
  const [commandOpen, setCommandOpen] = useState(false)

  useEffect(() => {
    const style = document.body.style
    const htmlStyle = document.documentElement.style

    style.overflow = "hidden"
    htmlStyle.overflow = "hidden"
    htmlStyle.overscrollBehavior = "none"
    htmlStyle.touchAction = "none"

    return () => {
      style.overflow = "unset"
      htmlStyle.overflow = "unset"
      htmlStyle.overscrollBehavior = "auto"
      htmlStyle.touchAction = "auto"
    }
  }, [])

  return (
    <div className="flex min-h-screen flex-col items-center justify-center overflow-hidden px-4 py-2 select-none">
      <motion.div
        initial={{ y: "100vh", opacity: 0 }}
        animate={{ y: -100, opacity: 1 }}
        className="select-none"
        transition={{ duration: 0.6, ease: "easeOut" }}
      >
        <Image
          src="/womm.png"
          alt="womm"
          draggable={false}
          width={250}
          className="transform rounded-full shadow-xl shadow-black/30 duration-250 hover:scale-110"
          height={250}
        />
      </motion.div>

      <motion.div
        initial={{ y: "100vh", opacity: 0 }}
        animate={{ y: -105, opacity: 1 }}
        transition={{ duration: 0.6, ease: "easeOut" }}
      >
        <p className="mt-3 text-2xl font-extrabold">
          Wo<span className="text-[#11e7e6]">mm</span>
        </p>
      </motion.div>

      <motion.div
        initial={{ y: "100vh", opacity: 0 }}
        animate={{ y: -110, opacity: 1 }}
        transition={{ duration: 0.6, ease: "easeOut" }}
      >
        <Button
          onClick={() => setCommandOpen(true)}
          variant="default"
          size="wide"
          className="mt-6 bg-[#11e7e6] px-6 py-2 font-extrabold text-white hover:bg-[#0db0b0]"
        >
          Get Started
        </Button>
      </motion.div>

      <motion.div
        initial={{ y: "100vh", opacity: 0 }}
        animate={{ y: -95, opacity: 1 }}
        transition={{ duration: 0.6, ease: "easeOut" }}
        className="mt-4 flex space-x-4 text-sm text-gray-400"
      >
        <div
          className="group flex cursor-pointer items-center space-x-2"
          onClick={() => open("https://github.com/mnsartawi/womm")}
        >
          <SiGithub size={16} />
          <span className="relative">
            GitHub
            <span className="absolute bottom-0 left-0 h-px w-0 bg-gray-400 transition-all duration-200 group-hover:w-full" />
          </span>
        </div>
        <p>;)</p>
        <div
          className="group flex cursor-pointer items-center space-x-1"
          onClick={() => open("https://womm.sartawi.dev")}
        >
          <Globe size={16} />
          <span className="relative">
            Website
            <span className="absolute bottom-0 left-0 h-px w-0 bg-gray-400 transition-all duration-200 group-hover:w-full" />
          </span>
        </div>
      </motion.div>

      <CommandDialog
        open={commandOpen}
        onOpenChange={setCommandOpen}
        className="sm:max-w-xl"
      >
        <Command className="[&_[cmdk-group]:not([hidden])_~[cmdk-group]]:pt-0 [&_[cmdk-input-wrapper]_svg]:h-5 [&_[cmdk-input-wrapper]_svg]:w-5 [&_[cmdk-item]_svg]:h-5 [&_[cmdk-item]_svg]:w-5 **:[[cmdk-group-heading]]:px-2 **:[[cmdk-group-heading]]:font-medium **:[[cmdk-group-heading]]:text-muted-foreground **:[[cmdk-group]]:px-2 **:[[cmdk-input]]:h-12 **:[[cmdk-item]]:px-2 **:[[cmdk-item]]:py-3">
          <CommandInput placeholder="Search pages..." />
          <CommandList>
            <CommandEmpty>No results found.</CommandEmpty>
            <CommandGroup heading="Navigation">
              <CommandItem redirect="/home">
                <FilePlus className="mr-2 size-6" />
                Generate *.womm
              </CommandItem>
              <CommandItem redirect="/home">
                <ScrollText className="mr-2 size-4" />
                Validate *.womm
              </CommandItem>
              <CommandItem redirect="/editor">
                <FileEdit className="mr-2 size-6" />
                Edit *.womm
              </CommandItem>
            </CommandGroup>
            <CommandGroup heading="Links">
              <CommandItem
                onSelect={() => {
                  setCommandOpen(false)
                  void open("https://github.com/mnsartawi/womm")
                }}
              >
                <SiGithub className="mr-2.5 size-5.5! shrink-0" />
                GitHub
              </CommandItem>
              <CommandItem
                onSelect={() => {
                  setCommandOpen(false)
                  router.push("/settings")
                }}
              >
                <CogIcon className="mr-2 size-5.5! shrink-0" />
                Settings
              </CommandItem>

              <CommandItem
                onSelect={() => {
                  setCommandOpen(false)
                  void open("https://womm.sartawi.dev")
                }}
              >
                <Image
                  className="mr-2 shrink-0 rounded-full object-contain"
                  src="/docs.webp"
                  draggable={false}
                  width={24}
                  height={24}
                  style={{ width: "24px", height: "24px", minWidth: "24px" }}
                  alt="docs"
                />
                Documentation
              </CommandItem>
              <CommandItem
                onSelect={() => {
                  setCommandOpen(false)
                  void open("https://womm.sartawi.dev")
                }}
              >
                <Image
                  className="mr-2 shrink-0 rounded-full object-contain"
                  src="/womm.png"
                  draggable={false}
                  width={24}
                  height={24}
                  style={{ width: "24px", height: "24px", minWidth: "24px" }}
                  alt="Website"
                />
                Website
              </CommandItem>
            </CommandGroup>
          </CommandList>
        </Command>
      </CommandDialog>
    </div>
  )
}
