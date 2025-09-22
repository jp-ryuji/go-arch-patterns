# Go ORM/Query Builder Selection Summary

## Background

Previously using SQLBoiler, which has entered **maintenance mode** with no active development. The need arose to find a modern, actively maintained alternative without losing SQLBoiler's strong type safety benefits (code generation from database schema).

## Requirements

- Type safety to minimize runtime errors
- PostgreSQL with pgx driver integration
- Efficient relationship loading (smart preloading to avoid both N+1 queries and massive JOIN result sets)

  ```go
  users := fetchUsers() // 1 query: SELECT * FROM users
  // ORM automatically does:
  books := fetchBooksByUserIDs(userIDs) // 1 query: SELECT * FROM books WHERE user_id IN (1,2,3...)
  // Then maps relationships in memory

  // This gives you:
  //
  // 2 total queries (not N+1)
  // Smaller result sets than massive JOINs
  // No duplicated data like JOINs would create
  ```

- SQL-like syntax preference
- Migration file auto-generation (optional but preferred)
- Long-term viability and community support
- **AI Assistant Support**: In the AI-assisted development era, popular libraries receive better support from AI assistants and have more comprehensive training data including Stack Overflow answers, tutorials, and examples

## Evaluated Options

### GORM (with gen)

**Pros:**

- Mature ecosystem with extensive documentation
- Excellent smart preloading with `Preload()` (avoids N+1 problem)
- Large community support
- Auto-migration support
- pgx compatibility

**Cons:**

- Less SQL-like syntax (method chaining)
- Can be "magical" and harder to debug
- Migration files not auto-generated from schema changes
- Performance overhead from reflection

### Ent (Meta/Facebook)

**Pros:**

- Excellent type safety through code generation
- Schema-first approach with automatic migration generation
- Strong smart preloading with graph traversal (avoids N+1 problem)
- Meta/Facebook backing ensures longevity
- Rich enterprise features (multi-tenancy, privacy policies)
- Growing community and documentation
- pgx compatibility

**Cons:**

- Less SQL-like syntax (graph-oriented)
- Steeper learning curve
- More opinionated approach
- Generated code can be heavy

### Bun

**Pros:**

- Very SQL-like syntax while maintaining reasonable type safety
- Good PostgreSQL/pgx integration
- Clean migration system
- Lightweight compared to full ORMs
- Efficient smart preloading with relations (avoids N+1 problem)

**Cons:**

- Weaker type safety (struct tags, runtime validation)
- Smaller community
- Manual schema definition
- Higher potential for runtime errors

### Bob

**Pros:**

- Generated from database schema (SQLBoiler successor)
- Better type safety than SQLBoiler
- SQL-like syntax
- Efficient smart preloading with ThenLoad (avoids N+1 problem)
- pgx compatibility

**Cons:**

- Very small community
- Limited ecosystem and documentation
- Less AI assistant support
- Newer project with uncertain long-term adoption

## Selection Criteria Analysis

| Criteria | GORM | Ent | Bun | Bob |
|----------|------|-----|-----|-----|
| Type Safety | ⭐⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐⭐ | ⭐⭐⭐⭐ |
| SQL-like Syntax | ⭐⭐ | ⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ |
| Relationship Loading | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐ | ⭐⭐⭐⭐ |
| pgx Integration | ⭐⭐⭐ | ⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐ |
| Auto Migrations | ⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐ | ⭐ |
| Community/Popularity | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐ | ⭐⭐⭐ | ⭐⭐ |
| Long-term Viability | ⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐⭐⭐ | ⭐⭐ |
| AI Assistant Support | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐ | ⭐⭐⭐ | ⭐⭐ |

## Decision Process

1. **Type Safety Priority**: Coming from SQLBoiler experience, runtime errors are a major concern. This eliminated Bun as the primary choice despite its excellent SQL-like syntax.

2. **Community & AI Era Considerations**: In the current AI-assisted development landscape, popularity and community support significantly impact development velocity and problem-solving capability.

3. **Long-term Viability**: Need for active maintenance and evolution, especially given SQLBoiler's maintenance-only status.

4. **PostgreSQL Focus**: All options support PostgreSQL well, but pgx integration varies in quality.

## Final Decision: Ent

**Selected: Ent** for the following reasons:

- **Type Safety**: Provides compile-time safety comparable to SQLBoiler through code generation
- **Automatic Migrations**: Schema-first approach with auto-generated migration files addresses our optional requirement perfectly
- **Enterprise Backing**: Meta/Facebook support ensures long-term viability and continued development
- **Community Growth**: Strong and growing community with excellent documentation
- **AI Assistant Support**: Better coverage in AI training data compared to smaller alternatives
- **PostgreSQL Integration**: Good pgx compatibility for performance benefits

**Trade-offs Accepted**:

- Less SQL-like syntax in favor of type safety and generated code benefits
- Steeper learning curve offset by comprehensive documentation and community support
- More opinionated approach provides consistency and prevents runtime errors
