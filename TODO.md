change add new products
<br/>
we need 3 page for this approach:
<br/>
1- add a product (images , videos, prices , title , slug ,sku , status )
<br/>
2- add a page for add attributes (product_attribute_values)
<br/>
3- add a page for add inventory for chosen attribute values

<br/>

explain WITH RECURSIVE CategoryHierarchy AS (
    SELECT id, title, parent_id
    FROM categories
    WHERE id = (SELECT category_id FROM products WHERE id = 1)

    UNION ALL

    SELECT c.id, c.title, c.parent_id
    FROM categories c
             INNER JOIN CategoryHierarchy ch ON c.id = ch.parent_id
)
SELECT *
FROM CategoryHierarchy
WHERE parent_id IS NULL
LIMIT 1;


<br/>

<ul>
<li>add a page for product inventory  </li>
<li>use transaction for add attr-values for a product</li>
<li>implement upload media for a product</li>


</ul>