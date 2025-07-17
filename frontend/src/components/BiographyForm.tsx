import React, { useState, useEffect } from 'react';
import { Control, useFieldArray, UseFormSetValue, UseFormWatch } from 'react-hook-form';
import { Plus, Trash2 } from 'lucide-react';
import { FormField } from './ui/FormField';
import Button from './ui/Button';
import { Card, CardHeader, CardContent } from './ui/Card';

interface ExperienceEntry {
  start_date: string;
  end_date: string;
  title: string;
  organization: string;
  description: string;
}

interface EducationEntry {
  start_date: string;
  end_date: string;
  title: string;
  institution: string;
}

interface Biography {
  experience: ExperienceEntry[];
  education: EducationEntry[];
}

interface BiographyFormProps {
  control: Control<any>;
  setValue: UseFormSetValue<any>;
  watch: UseFormWatch<any>;
  errors?: any;
}

const BiographyForm: React.FC<BiographyFormProps> = ({
  control,
  setValue,
  watch,
  errors,
}) => {
  const {
    fields: experienceFields,
    append: appendExperience,
    remove: removeExperience,
  } = useFieldArray({
    control,
    name: 'biography.experience',
  });

  const {
    fields: educationFields,
    append: appendEducation,
    remove: removeEducation,
  } = useFieldArray({
    control,
    name: 'biography.education',
  });

  const addExperienceEntry = () => {
    appendExperience({
      start_date: '',
      end_date: '',
      title: '',
      organization: '',
      description: '',
    });
  };

  const addEducationEntry = () => {
    appendEducation({
      start_date: '',
      end_date: '',
      title: '',
      institution: '',
    });
  };

  // Initialize with one entry of each type if none exist
  useEffect(() => {
    if (experienceFields.length === 0) {
      addExperienceEntry();
    }
    if (educationFields.length === 0) {
      addEducationEntry();
    }
  }, []);

  return (
    <div className="space-y-6">
      {/* Experience Section */}
      <Card>
        <CardHeader>
          <div className="flex justify-between items-center">
            <h3 className="text-lg font-semibold">Professional Experience</h3>
            <Button
              type="button"
              variant="outline"
              size="sm"
              onClick={addExperienceEntry}
              className="flex items-center gap-2"
            >
              <Plus className="h-4 w-4" />
              Add Experience
            </Button>
          </div>
        </CardHeader>
        <CardContent>
          <div className="space-y-4">
            {experienceFields.map((field, index) => (
              <div key={field.id} className="border rounded-lg p-4 space-y-4">
                <div className="flex justify-between items-center">
                  <h4 className="font-medium">Experience {index + 1}</h4>
                  {experienceFields.length > 1 && (
                    <Button
                      type="button"
                      variant="ghost"
                      size="sm"
                      onClick={() => removeExperience(index)}
                      className="text-red-600 hover:text-red-800"
                    >
                      <Trash2 className="h-4 w-4" />
                    </Button>
                  )}
                </div>
                
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <FormField
                    name={`biography.experience.${index}.start_date`}
                    label="Start Date"
                    type="month"
                    control={control}
                    error={errors?.biography?.experience?.[index]?.start_date}
                    required
                  />
                  <FormField
                    name={`biography.experience.${index}.end_date`}
                    label="End Date"
                    type="month"
                    control={control}
                    error={errors?.biography?.experience?.[index]?.end_date}
                    required
                    placeholder="YYYY-MM or 'Present'"
                  />
                </div>
                
                <FormField
                  name={`biography.experience.${index}.title`}
                  label="Job Title"
                  control={control}
                  error={errors?.biography?.experience?.[index]?.title}
                  placeholder="e.g., Senior Software Engineer"
                  required
                />
                
                <FormField
                  name={`biography.experience.${index}.organization`}
                  label="Organization"
                  control={control}
                  error={errors?.biography?.experience?.[index]?.organization}
                  placeholder="e.g., Company Name"
                  required
                />
                
                <FormField
                  name={`biography.experience.${index}.description`}
                  label="Description"
                  type="textarea"
                  control={control}
                  error={errors?.biography?.experience?.[index]?.description}
                  placeholder="Describe your responsibilities and achievements..."
                  required
                  rows={3}
                />
              </div>
            ))}
          </div>
        </CardContent>
      </Card>

      {/* Education Section */}
      <Card>
        <CardHeader>
          <div className="flex justify-between items-center">
            <h3 className="text-lg font-semibold">Education</h3>
            <Button
              type="button"
              variant="outline"
              size="sm"
              onClick={addEducationEntry}
              className="flex items-center gap-2"
            >
              <Plus className="h-4 w-4" />
              Add Education
            </Button>
          </div>
        </CardHeader>
        <CardContent>
          <div className="space-y-4">
            {educationFields.map((field, index) => (
              <div key={field.id} className="border rounded-lg p-4 space-y-4">
                <div className="flex justify-between items-center">
                  <h4 className="font-medium">Education {index + 1}</h4>
                  {educationFields.length > 1 && (
                    <Button
                      type="button"
                      variant="ghost"
                      size="sm"
                      onClick={() => removeEducation(index)}
                      className="text-red-600 hover:text-red-800"
                    >
                      <Trash2 className="h-4 w-4" />
                    </Button>
                  )}
                </div>
                
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <FormField
                    name={`biography.education.${index}.start_date`}
                    label="Start Date"
                    type="month"
                    control={control}
                    error={errors?.biography?.education?.[index]?.start_date}
                    required
                  />
                  <FormField
                    name={`biography.education.${index}.end_date`}
                    label="End Date"
                    type="month"
                    control={control}
                    error={errors?.biography?.education?.[index]?.end_date}
                    required
                  />
                </div>
                
                <FormField
                  name={`biography.education.${index}.title`}
                  label="Degree/Qualification"
                  control={control}
                  error={errors?.biography?.education?.[index]?.title}
                  placeholder="e.g., Master of Science in Computer Science"
                  required
                />
                
                <FormField
                  name={`biography.education.${index}.institution`}
                  label="Institution"
                  control={control}
                  error={errors?.biography?.education?.[index]?.institution}
                  placeholder="e.g., University Name"
                  required
                />
              </div>
            ))}
          </div>
        </CardContent>
      </Card>

      {/* Biography Preview */}
      <Card>
        <CardHeader>
          <h3 className="text-lg font-semibold">Biography Preview</h3>
        </CardHeader>
        <CardContent>
          <BiographyPreview watch={watch} />
        </CardContent>
      </Card>
    </div>
  );
};

// Biography Preview Component
interface BiographyPreviewProps {
  watch: UseFormWatch<any>;
}

const BiographyPreview: React.FC<BiographyPreviewProps> = ({ watch }) => {
  const biography = watch('biography');
  
  if (!biography || (!biography.experience?.length && !biography.education?.length)) {
    return <p className="text-gray-500 italic">Biography preview will appear here as you fill in the form.</p>;
  }

  return (
    <div className="space-y-4">
      {/* Education Section */}
      {biography.education?.length > 0 && (
        <div>
          <h4 className="font-semibold text-base mb-2">Education</h4>
          <ul className="space-y-1">
            {biography.education.map((edu: EducationEntry, index: number) => (
              <li key={index} className="text-sm">
                • {edu.start_date} - {edu.end_date} - {edu.title} - {edu.institution}
              </li>
            ))}
          </ul>
        </div>
      )}

      {/* Experience Section */}
      {biography.experience?.length > 0 && (
        <div>
          <h4 className="font-semibold text-base mb-2">Experience</h4>
          <ul className="space-y-1">
            {biography.experience.map((exp: ExperienceEntry, index: number) => (
              <li key={index} className="text-sm">
                • {exp.start_date} - {exp.end_date} - {exp.title} - {exp.organization} - {exp.description}
              </li>
            ))}
          </ul>
        </div>
      )}
    </div>
  );
};

export default BiographyForm;