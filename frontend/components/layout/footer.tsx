import Link from 'next/link';
import Image from 'next/image';

export function Footer() {
  return (
    <footer className="border-t bg-background">
      <div className="container flex flex-col items-center justify-between gap-4 py-10 md:h-24 md:flex-row md:py-0">
        <div className="flex flex-col items-center gap-4 px-8 md:flex-row md:gap-4 md:px-0">
          <Image 
            src="/images/logo/Icon Logo - Color.svg"
            alt="BQA Logo"
            width={32}
            height={32}
            className="h-8 w-auto"
          />
          <p className="text-center text-sm leading-loose text-muted-foreground md:text-left">
            © {new Date().getFullYear()} <span className="text-primary font-medium">BQA ExpertDB</span>. All rights reserved.
          </p>
        </div>
        <div className="flex gap-4">
          <Link href="/about" className="text-sm font-medium text-primary hover:underline underline-offset-4">
            About
          </Link>
          <Link href="/contact" className="text-sm font-medium text-primary hover:underline underline-offset-4">
            Contact
          </Link>
        </div>
      </div>
    </footer>
  );
}