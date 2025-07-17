import React from 'react';

interface CheckboxProps {
  id: string;
  checked: boolean;
  disabled?: boolean;
  onChange: () => void;
  label?: string;
  className?: string;
}

export const Checkbox: React.FC<CheckboxProps> = ({
  id,
  checked,
  disabled = false,
  onChange,
  label,
  className = ''
}) => {
  return (
    <div className={`flex items-center ${className}`}>
      <input
        type="checkbox"
        id={id}
        checked={checked}
        disabled={disabled}
        onChange={onChange}
        className="h-4 w-4 text-primary focus:ring-primary border-neutral-300 rounded disabled:opacity-50"
      />
      {label && (
        <label
          htmlFor={id}
          className={`ml-2 text-sm ${
            disabled ? 'text-neutral-400' : 'text-neutral-700'
          } cursor-pointer`}
        >
          {label}
        </label>
      )}
    </div>
  );
};