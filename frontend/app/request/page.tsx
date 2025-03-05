'use client';

import { RequestForm } from './request-form';
import { Navbar } from '@/components/layout/navbar';
import RequireAuth from '@/components/auth/require-auth';

export default function RequestPage() {
  return (
    <RequireAuth>
      <>
        <Navbar />
        <div className="container py-10">
          <div className="mx-auto max-w-3xl">
            <h1 className="text-3xl font-bold mb-6">Submit Expert Request</h1>
            <p className="text-muted-foreground mb-8">
              Please fill out the form below to submit a new expert for review. Our team will review the information and add the expert to our database if approved.
            </p>
            <RequestForm />
          </div>
        </div>
      </>
    </RequireAuth>
  );
}