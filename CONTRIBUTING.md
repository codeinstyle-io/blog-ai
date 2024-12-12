# Contributing to Captain

Welcome to Captain! We're excited that you want to contribute. This document outlines our contribution process and guidelines.

## AI-First Development

Captain is an AI-first project, which means we actively encourage the use of AI tools in development. We believe AI can help create better, more maintainable code and improve developer productivity.

- Feel free to use AI assistants (like Claude, GPT-4, etc.) for:
  - Code generation and refactoring
  - Documentation writing
  - Test case creation
  - Bug fixing
  - Architecture discussions
- When submitting PRs, please document which AI tools you used (see PR template)

## Branch Naming Convention

We use a descriptive branch naming scheme to clearly communicate the purpose of each branch:

- For new features: `feat.my-feature-name`
  - Example: `feat.user-authentication`
  - Example: `feat.markdown-support`

- For bug fixes: `fix.bug-description`
  - Example: `fix.login-validation`
  - Example: `fix.memory-leak`

- For chores and maintenance: `chore.task-description`
  - Example: `chore.updated-documentation`
  - Example: `chore.dependency-upgrade`

## Development Workflow

1. Fork the repository and create your branch from `main`
2. Follow the branch naming convention described above
3. Make your changes:
   - Write meaningful commit messages
   - Add tests if applicable
   - Update documentation as needed
4. Ensure all tests pass by running `make test`
5. Create a Pull Request

## Pull Request Process

1. Use our PR template (located at `.github/pull_request_template.md`)
2. Fill in all required sections, including:
   - AI tools used
   - Editor/IDE used
   - Prompts used (if applicable)
3. Link any related issues
4. Update the README.md if your changes affect it
5. Wait for review from maintainers

## Code Style

- Follow existing code style and patterns
- Use meaningful variable and function names
- Write clear comments for complex logic
- Keep functions focused and concise
- Add appropriate error handling

## Testing

- Write tests for new features
- Update tests for modified code
- Ensure all tests pass before submitting PR
- Aim for good test coverage

## Documentation

- Update documentation for new features
- Keep documentation clear and concise
- Include code examples where helpful
- Check for spelling and grammar

## Need Help?

- Check existing issues and documentation
- Create a new issue for questions
- Join our community discussions

## License

By contributing, you agree that your contributions will be licensed under the project's license.

Thank you for contributing to Captain! 🚀
