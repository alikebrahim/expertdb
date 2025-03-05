import Link from 'next/link';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardFooter } from '@/components/ui/card';

export const metadata = {
  title: 'Request Submitted - ExpertDB',
  description: 'Your expert request has been successfully submitted',
};

export default function RequestSuccessPage() {
  return (
    <div className="container py-10">
      <div className="mx-auto max-w-md">
        <Card className="border-green-200">
          <CardContent className="pt-6 text-center">
            <div className="mb-4 rounded-full bg-green-50 p-3 inline-flex mx-auto">
              <svg 
                xmlns="http://www.w3.org/2000/svg" 
                className="h-8 w-8 text-green-600" 
                viewBox="0 0 24 24" 
                fill="none" 
                stroke="currentColor" 
                strokeWidth="2" 
                strokeLinecap="round" 
                strokeLinejoin="round"
              >
                <path d="M22 11.08V12a10 10 0 1 1-5.93-9.14" />
                <polyline points="22 4 12 14.01 9 11.01" />
              </svg>
            </div>
            <h2 className="text-2xl font-bold mb-2">Request Submitted</h2>
            <p className="text-muted-foreground mb-4">
              Your expert request has been successfully submitted for review. Our team will process your request and you will be notified once it has been approved.
            </p>
          </CardContent>
          <CardFooter className="flex flex-col space-y-2">
            <Button asChild className="w-full">
              <Link href="/">Return to Home</Link>
            </Button>
            <Button asChild variant="outline" className="w-full">
              <Link href="/search">Search Experts</Link>
            </Button>
          </CardFooter>
        </Card>
      </div>
    </div>
  );
}