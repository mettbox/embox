import { Ref, reactive, computed } from 'vue';

/**
 * FormField interface represents a form field with its validation rules.
 *
 * @interface FormField
 * @property {Ref<any>} ref - A reference to the form field element.
 * @property {() => boolean} [shouldValidate] - A function that determines if the field should be validated.
 * @property {Array<{ validate: (value: string) => boolean; message: string }>} rules - An array of validation rules for the field.
 */
export interface FormField {
  ref: Ref<any>;
  shouldValidate?: () => boolean;
  rules: {
    validate: (value: string) => boolean;
    message: string;
  }[];
}

/**
 * Validates a form field based on the provided rules and model.
 *
 * @param {Record<string, FormField>} fields - Fields to validate
 * @param {Record<string, any>} model - Model to validate against
 * @returns
 */
export function useValidation(fields: Record<string, FormField>, model: Record<string, any>) {
  const fieldErrors = reactive<Record<string, string>>({});
  const isTouched = reactive<Record<string, boolean>>({});

  function validateField(name: string): void {
    const field = fields[name];
    const el = field.ref.value?.$el;
    if (!el) return;

    if (field.shouldValidate && !field.shouldValidate()) {
      fieldErrors[name] = '';
      isTouched[name] = false;
      return;
    }

    const value = model[name];

    let errorMessage = '';
    for (const rule of field.rules) {
      if (!rule.validate(value)) {
        errorMessage = rule.message;
        break;
      }
    }

    fieldErrors[name] = errorMessage;
    isTouched[name] = true;

    el.classList.remove('ion-valid', 'ion-invalid');
    el.classList.add('ion-touched');

    if (errorMessage) {
      el.classList.add('ion-invalid');
    } else {
      el.classList.add('ion-valid');
    }
  }

  function validateAllFields(): void {
    Object.keys(fields).forEach(validateField);
  }

  const isFormValid = computed(
    // () => Object.keys(fields).every((key) => isTouched[key]) && Object.values(fieldErrors).every((e) => e === ''),
    () => Object.values(fieldErrors).every((e) => e === ''),
  );

  return {
    validateField,
    validateAllFields,
    isFormValid,
    fieldErrors,
    isTouched,
  };
}
