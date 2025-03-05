import * as React from "react"
import { Slot } from "@radix-ui/react-slot"
import { cva, type VariantProps } from "class-variance-authority"

import { cn } from "@/lib/utils"

const buttonVariants = cva(
  "inline-flex items-center justify-center whitespace-nowrap rounded-[4px] text-sm font-medium transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50",
  {
    variants: {
      variant: {
        default: "bg-[#133566] text-white hover:bg-[#1B4882]",
        destructive:
          "bg-[#FF4040] text-white hover:bg-[#E03131]",
        outline:
          "border border-[#133566] bg-white text-[#133566] hover:bg-[#F5F5F5] hover:text-[#133566]",
        outlineInverse:
          "border border-white bg-transparent text-white hover:bg-[#1B4882] hover:text-white hover:border-white",
        secondary:
          "bg-[#1B4882] text-white hover:bg-[#1B4882]/80",
        success:
          "bg-[#192012] text-white hover:bg-[#192012]/90",
        ghost: "hover:bg-[#DC8335]/10 hover:text-[#DC8335]",
        link: "text-[#133566] underline-offset-4 hover:underline",
      },
      size: {
        default: "h-10 px-4 py-2",
        sm: "h-9 px-3",
        lg: "h-11 px-8",
        icon: "h-10 w-10",
      },
    },
    defaultVariants: {
      variant: "default",
      size: "default",
    },
  }
)

export interface ButtonProps
  extends React.ButtonHTMLAttributes<HTMLButtonElement>,
    VariantProps<typeof buttonVariants> {
  asChild?: boolean
}

const Button = React.forwardRef<HTMLButtonElement, ButtonProps>(
  ({ className, variant, size, asChild = false, ...props }, ref) => {
    const Comp = asChild ? Slot : "button"
    return (
      <Comp
        className={cn(buttonVariants({ variant, size, className }))}
        ref={ref}
        {...props}
      />
    )
  }
)
Button.displayName = "Button"

export { Button, buttonVariants }