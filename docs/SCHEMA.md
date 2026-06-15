# Database schema

## users

Represents an account on the platform.

| Column         | Type           | Constraints             | Notes                                              |
| -------------- | -------------- | ----------------------- | -------------------------------------------------- |
| `user_id`      | `bigint`       | PK                      | Auto-incremented surrogate key                     |
| `handle`       | `varchar(50)`  | UK, NOT NULL            | e.g. `@person` — used in mentions and route params |
| `display_name` | `varchar(100)` | NOT NULL                | Shown in the UI                                    |
| `biography`    | `text`         |                         | Profile description                                |
| `created_at`   | `timestamptz`  | NOT NULL, DEFAULT now() |                                                    |
| `updated_at`   | `timestamptz`  | DEFAULT now()           |                                                    |
| `deleted_at`   | `timestamptz`  | DEFAULT now()           |                                                    |

---

## posts

A flat post authored by a user. No nesting or media.

| Column       | Type          | Constraints             | Notes                              |
| ------------ | ------------- | ----------------------- | ---------------------------------- |
| `post_id`    | `bigint`      | PK                      | Auto-incremented surrogate key     |
| `author_id`  | `bigint`      | FK → users.id, NOT NULL | Cascade delete on user removal     |
| `body`       | `text`        | NOT NULL                | Post content                       |
| `like_count` | `int`         | NOT NULL, DEFAULT 0     | Denormalized — maintained on write |
| `created_at` | `timestamptz` | NOT NULL, DEFAULT now() |                                    |
| `updated_at` | `timestamptz` | DEFAULT now()           |                                    |
| `deleted_at` | `timestamptz` | DEFAULT now()           |                                    |

---

## follows

A directional relationship between two users.

| Column        | Type          | Constraints                   | Notes                     |
| ------------- | ------------- | ----------------------------- | ------------------------- |
| `follower_id` | `bigint`      | PK (composite), FK → users.id | The user who is following |
| `followee_id` | `bigint`      | PK (composite), FK → users.id | The user being followed   |
| `created_at`  | `timestamptz` | NOT NULL, DEFAULT now()       |                           |

**Composite PK** on `(follower_id, followee_id)` enforces uniqueness and prevents duplicate follows.

---

## likes

A user liking a post. One row per unique user–post pair.

| Column       | Type          | Constraints                   | Notes |
| ------------ | ------------- | ----------------------------- | ----- |
| `user_id`    | `bigint`      | PK (composite), FK → users.id |       |
| `post_id`    | `bigint`      | PK (composite), FK → posts.id |       |
| `created_at` | `timestamptz` | NOT NULL, DEFAULT now()       |       |

**Composite PK** on `(user_id, post_id)` makes likes idempotent — re-inserting has no effect.

---

## Key indexes

| Index                      | Table     | Columns                        | Purpose                                       |
| -------------------------- | --------- | ------------------------------ | --------------------------------------------- |
| `idx_posts_author_created` | `posts`   | `(author_id, created_at DESC)` | Feed query — fetch posts by followed accounts |
| `idx_follows_follower`     | `follows` | `(follower_id)`                | Look up who a user follows                    |
| `idx_follows_followee`     | `follows` | `(followee_id)`                | Look up a user's followers                    |
| `idx_likes_post`           | `likes`   | `(post_id)`                    | Count or list likes per post                  |

---
