import React, { useState, useEffect } from 'react';

interface Step {
  id: string | number;
  label: string;
  description?: string;
  optional?: boolean;
}

interface ProgressStepperProps {
  steps: Step[];
  currentStep: number;
  onStepClick?: (step: number) => void;
  orientation?: 'horizontal' | 'vertical';
  size?: 'sm' | 'md' | 'lg';
  className?: string;
  animated?: boolean;
  showPercentage?: boolean;
  allowNavigation?: boolean;
}

export const ProgressStepper: React.FC<ProgressStepperProps> = ({
  steps,
  currentStep,
  onStepClick,
  orientation = 'horizontal',
  size = 'md',
  className = '',
  animated = true,
  showPercentage = false,
  allowNavigation = true,
}) => {
  const [progress, setProgress] = useState(0);
  const totalSteps = steps.length;

  useEffect(() => {
    const percentage = ((currentStep - 1) / (totalSteps - 1)) * 100;
    setProgress(Math.max(0, Math.min(100, percentage)));
  }, [currentStep, totalSteps]);

  // Size classes
  const sizeClasses = {
    sm: {
      stepper: 'text-xs',
      indicator: 'h-6 w-6 text-xs',
      line: 'h-0.5',
    },
    md: {
      stepper: 'text-sm',
      indicator: 'h-8 w-8 text-sm',
      line: 'h-1',
    },
    lg: {
      stepper: 'text-base',
      indicator: 'h-10 w-10 text-base',
      line: 'h-1.5',
    },
  };

  const selectedSize = sizeClasses[size];

  // Determine if horizontal or vertical
  const isHorizontal = orientation === 'horizontal';

  const handleStepClick = (stepIndex: number) => {
    if (allowNavigation && stepIndex <= currentStep && onStepClick) {
      onStepClick(stepIndex);
    }
  };

  return (
    <div 
      className={`${isHorizontal ? 'w-full' : 'h-full'} ${selectedSize.stepper} ${className}`}
      aria-label="Progress"
    >
      <div 
        className={`flex ${isHorizontal ? 'flex-row' : 'flex-col'} ${isHorizontal ? 'w-full' : 'h-full'}`}
      >
        {steps.map((step, index) => {
          const isCompleted = index < currentStep;
          const isCurrent = index === currentStep;
          const isClickable = allowNavigation && index <= currentStep;

          return (
            <div 
              key={step.id} 
              className={`${isHorizontal ? 'flex-1' : 'flex-none'} relative ${isHorizontal ? '' : 'pb-8'}`}
            >
              {/* Connector Line */}
              {index > 0 && (
                <div 
                  className={`
                    absolute 
                    ${isHorizontal ? 'top-1/2 -translate-y-1/2 left-0' : 'top-0 left-1/2 -translate-x-1/2'} 
                    ${isHorizontal ? 'h-0.5 bg-gray-200 dark:bg-gray-700' : 'w-0.5 bg-gray-200 dark:bg-gray-700'} 
                    ${isHorizontal ? 'w-full -left-full' : 'h-full'}
                  `}
                >
                  <div 
                    className={`
                      ${isHorizontal ? 'h-full' : 'w-full'} 
                      bg-blue-600 
                      transition-all duration-500 ease-in-out
                    `}
                    style={{ 
                      [isHorizontal ? 'width' : 'height']: index <= currentStep ? '100%' : '0%'
                    }}
                  />
                </div>
              )}

              {/* Step Indicator + Label */}
              <div className={`${isHorizontal ? '' : 'flex flex-row items-center'} group relative`}>
                {/* Step Number/Icon */}
                <div
                  onClick={() => handleStepClick(index)}
                  className={`
                    ${selectedSize.indicator}
                    rounded-full
                    flex items-center justify-center
                    ${isCompleted ? 'bg-blue-600 text-white' : isCurrent ? 'bg-blue-100 text-blue-600 border-2 border-blue-600' : 'bg-gray-100 text-gray-500 dark:bg-gray-700 dark:text-gray-400'}
                    ${isClickable ? 'cursor-pointer hover:ring-2 hover:ring-offset-2 hover:ring-blue-500' : ''}
                    transition-all duration-300
                    z-10
                  `}
                  aria-current={isCurrent ? 'step' : undefined}
                >
                  {isCompleted ? (
                    <svg 
                      className="w-5 h-5" 
                      fill="currentColor" 
                      viewBox="0 0 20 20" 
                      xmlns="http://www.w3.org/2000/svg"
                    >
                      <path 
                        fillRule="evenodd" 
                        d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" 
                        clipRule="evenodd" 
                      />
                    </svg>
                  ) : (
                    <span>{index + 1}</span>
                  )}
                </div>

                {/* Step Label */}
                <div 
                  className={`
                    ${isHorizontal ? 'mt-2' : 'ml-3'}
                    ${isCompleted || isCurrent ? 'text-gray-900 dark:text-gray-100' : 'text-gray-500 dark:text-gray-400'}
                  `}
                >
                  <span 
                    className={`
                      font-medium 
                      ${isClickable ? 'cursor-pointer' : ''}
                    `}
                    onClick={() => handleStepClick(index)}
                  >
                    {step.label} {step.optional && <span className="text-gray-400 dark:text-gray-500 text-xs">(Optional)</span>}
                  </span>
                  {step.description && (
                    <p className="text-gray-500 dark:text-gray-400 text-xs">
                      {step.description}
                    </p>
                  )}
                </div>
              </div>
            </div>
          );
        })}
      </div>

      {/* Progress percentage indicator */}
      {showPercentage && (
        <div className="mt-4 text-right text-sm text-gray-500 dark:text-gray-400">
          {Math.round(progress)}% Complete
        </div>
      )}
    </div>
  );
};