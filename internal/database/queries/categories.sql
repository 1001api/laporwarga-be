-- name: CreateCategory :one
INSERT INTO categories (
    name, 
    slug,
    icon,
    color,
    is_active,
    sort_order
) VALUES (
    @name,
    @slug,
    @icon,
    @color,
    @is_active,
    @sort_order
) RETURNING id;
       
-- name: CheckCategoryExist :one
SELECT EXISTS (
    SELECT 1 
    FROM categories 
    WHERE name = @name OR slug = @slug
);

-- name: GetCategories :many
SELECT 
    id,
    name,
    slug,
    icon,
    color,
    is_active,
    sort_order,
    created_at,
    updated_at
FROM categories
WHERE is_active = TRUE AND deleted_at IS NULL
ORDER BY sort_order ASC;

-- name: GetCategoryById :one
SELECT 
    id,
    name,
    slug,
    icon,
    color,
    is_active,
    sort_order,
    created_at,
    updated_at
FROM categories
WHERE id = @id AND deleted_at IS NULL;

-- name: GetCategoryBySlug :one
SELECT 
    id,
    name,
    slug,
    icon,
    color,
    is_active,
    sort_order,
    created_at,
    updated_at
FROM categories
WHERE slug = @slug AND deleted_at IS NULL;

-- name: SearchCategories :many
SELECT 
    id,
    name,
    slug,
    icon,
    color,
    is_active,
    sort_order,
    created_at,
    updated_at
FROM categories
WHERE name ILIKE '%' || @search_term::text || '%' AND deleted_at IS NULL
ORDER BY 
    CASE WHEN @sort_by::text = 'name'       AND @sort_order::text = 'ASC'  THEN name       END ASC,
    CASE WHEN @sort_by::text = 'name'       AND @sort_order::text = 'DESC' THEN name       END DESC,
    CASE WHEN @sort_by::text = 'slug'       AND @sort_order::text = 'ASC'  THEN slug       END ASC,
    CASE WHEN @sort_by::text = 'slug'       AND @sort_order::text = 'DESC' THEN slug       END DESC,
    CASE WHEN @sort_by::text = 'icon'       AND @sort_order::text = 'ASC'  THEN icon       END ASC,
    CASE WHEN @sort_by::text = 'icon'       AND @sort_order::text = 'DESC' THEN icon       END DESC,
    CASE WHEN @sort_by::text = 'color'      AND @sort_order::text = 'ASC'  THEN color      END ASC,
    CASE WHEN @sort_by::text = 'color'      AND @sort_order::text = 'DESC' THEN color      END DESC,
    CASE WHEN @sort_by::text = 'is_active'  AND @sort_order::text = 'ASC'  THEN is_active  END ASC,
    CASE WHEN @sort_by::text = 'is_active'  AND @sort_order::text = 'DESC' THEN is_active  END DESC,
    CASE WHEN @sort_by::text = 'sort_order' AND @sort_order::text = 'ASC'  THEN sort_order END ASC,
    CASE WHEN @sort_by::text = 'sort_order' AND @sort_order::text = 'DESC' THEN sort_order END DESC;
        

-- name: ToggleCategoryActiveStatus :one
UPDATE categories
SET is_active = NOT is_active
WHERE id = @id AND deleted_at IS NULL 
RETURNING id, is_active;

-- name: UpdateCategory :one
UPDATE categories
SET 
    name = CASE
        WHEN @name::text IS NOT NULL AND @name::text != name THEN @name::text
        ELSE name
    END,
    slug = CASE
        WHEN @slug::text IS NOT NULL AND @slug::text != slug THEN @slug::text
        ELSE slug
    END,
    icon = CASE
        WHEN @icon::text IS NOT NULL AND @icon::text != icon THEN @icon::text
        ELSE icon
    END,
    color = CASE
        WHEN @color::text IS NOT NULL AND @color::text != color THEN @color::text
        ELSE color
    END,
    is_active = CASE
        WHEN @is_active::boolean IS NOT NULL AND @is_active::boolean != is_active THEN @is_active::boolean
        ELSE is_active
    END,
    sort_order = CASE
        WHEN @sort_order::integer IS NOT NULL AND @sort_order::integer != sort_order THEN @sort_order::integer
        ELSE sort_order
    END
WHERE id = @id AND deleted_at IS NULL 
RETURNING id;

-- name: DeleteCategory :one
UPDATE categories
SET deleted_at = NOW()
WHERE id = @id AND deleted_at IS NULL 
RETURNING id;