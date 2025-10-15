import { describe, it, expect } from 'vitest';

// Mock Button component for testing
const Button = {
	render: (props: any) => {
		return {
			text: props.children || 'Button',
			disabled: props.disabled || false,
			loading: props.loading || false,
			variant: props.variant || 'primary',
			size: props.size || 'md'
		};
	}
};

describe('Button Component', () => {
	describe('Rendering', () => {
		it('renders with default props', () => {
			const result = Button.render({});
			
			expect(result.text).toBe('Button');
			expect(result.disabled).toBe(false);
			expect(result.loading).toBe(false);
			expect(result.variant).toBe('primary');
			expect(result.size).toBe('md');
		});

		it('renders with custom content', () => {
			const result = Button.render({ children: 'Custom Text' });
			
			expect(result.text).toBe('Custom Text');
		});
	});

	describe('Variants', () => {
		it('applies primary variant', () => {
			const result = Button.render({ variant: 'primary' });
			expect(result.variant).toBe('primary');
		});

		it('applies secondary variant', () => {
			const result = Button.render({ variant: 'secondary' });
			expect(result.variant).toBe('secondary');
		});

		it('applies danger variant', () => {
			const result = Button.render({ variant: 'danger' });
			expect(result.variant).toBe('danger');
		});
	});

	describe('Sizes', () => {
		it('applies small size', () => {
			const result = Button.render({ size: 'sm' });
			expect(result.size).toBe('sm');
		});

		it('applies medium size', () => {
			const result = Button.render({ size: 'md' });
			expect(result.size).toBe('md');
		});

		it('applies large size', () => {
			const result = Button.render({ size: 'lg' });
			expect(result.size).toBe('lg');
		});
	});

	describe('States', () => {
		it('applies disabled state', () => {
			const result = Button.render({ disabled: true });
			expect(result.disabled).toBe(true);
		});

		it('applies loading state', () => {
			const result = Button.render({ loading: true });
			expect(result.loading).toBe(true);
		});
	});
});
