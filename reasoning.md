REASONING & DESIGN DECISIONS

1.ARCHITECTURE CHOICE - Used layered architecture: handler → service →
repository - Initially mixed validation and business logic - Debugging
became difficult due to tight coupling - Refactoring improved
readability and maintainability

2.DYNAMIC AGE CALCULATION - Did not store age in the database - Age is
derived data and changes over time - Calculated dynamically using Go’s
time package - Prevents inconsistency and avoids extra updates

3.VALIDATION APPROACH - Used go-playground/validator - Added required
field checks - Added name length validation - Added strict date format
validation - Faced integration issues initially - Solved by isolating
validation into an internal helper

4.EDGE CASE HANDLING - Empty names - Invalid date formats - Future DOB
values - Non-existent users - Consistent and readable error responses

5.MIDDLEWARE & LOGGING - Injected requestId - Logged request duration -
Improved traceability - Used Uber Zap for structured logging

6.TESTING FOCUS - Unit-tested age calculation - Covered birthday edge
cases - Avoided off-by-one errors

WHAT I LEARNED - Separation of concerns reduces debugging time -
Validation needs structure - Middleware improves observability -
Refactoring improves clarity