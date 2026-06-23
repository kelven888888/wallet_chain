INSERT INTO `nov_role` (`id`, `pid`, `name`, `icon`, `is_menu`, `desc`, `module`, `action`, `sort`, `is_default`, `created_at`, `updated_at`, `deleted_at`) VALUES (1013, 54, '商品类别管理', '', 1, NULL, 'Category', 'Index', 0, 0, '2024-07-13 12:51:12.094', '2024-06-09 12:03:11.987', NULL);
INSERT INTO `nov_role` (`id`, `pid`, `name`, `icon`, `is_menu`, `desc`, `module`, `action`, `sort`, `is_default`, `created_at`, `updated_at`, `deleted_at`) VALUES (NULL, 1013, '商品类别编辑', '', 0, NULL, 'Category', 'Edit', 0, 0, '2024-07-13 12:51:12.112', '2024-06-09 12:03:11.987', NULL);
INSERT INTO `nov_role` (`id`, `pid`, `name`, `icon`, `is_menu`, `desc`, `module`, `action`, `sort`, `is_default`, `created_at`, `updated_at`, `deleted_at`) VALUES (NULL, 1013, '商品类别删除', '', 0, NULL, 'Category', 'Delete', 0, 0, '2024-07-13 12:51:12.114', '2024-06-09 12:03:11.987', NULL);
INSERT INTO `nov_role` (`id`, `pid`, `name`, `icon`, `is_menu`, `desc`, `module`, `action`, `sort`, `is_default`, `created_at`, `updated_at`, `deleted_at`) VALUES (NULL, 1013, '商品类别添加', '', 0, NULL, 'Category', 'Add', 0, 0, '2024-07-13 12:51:12.115', '2024-06-09 12:03:11.987', NULL);
INSERT INTO `nov_role` (`id`, `pid`, `name`, `icon`, `is_menu`, `desc`, `module`, `action`, `sort`, `is_default`, `created_at`, `updated_at`, `deleted_at`) VALUES (NULL, 1013, '商品类别批量删除', '', 0, NULL, 'Category', 'Deletebatch', 0, 0, '2024-07-13 12:51:12.117', '2024-06-09 12:03:11.987', NULL);
package model

import "time" 
type Category struct {
Model model.Model `comment:"" types:"" text:"" json:"model" form:"model" range:"" edit:""  `
Name string `comment:"分类名称" types:"" text:"" json:"name" form:"name" range:"" edit:""  `
Pid int64 `comment:"" types:"" text:"" json:"pid" form:"pid" range:"" edit:""  `
Icon string `comment:"" types:"" text:"" json:"icon" form:"icon" range:"" edit:"0"  `
State int64 `comment:"分类状态" types:"radio" text:"未启用,已启用" json:"state" form:"state" range:"0,1" edit:""  `
Sort int64 `comment:"排序" types:"" text:"" json:"sort" form:"sort" range:"" edit:""  `
Tag string `comment:"" types:"" text:"" json:"tag" form:"tag" range:"" edit:"0"  `
Children *model.CategoryTrees `comment:"" types:"" text:"" json:"children" form:"children" range:"" edit:""  `
}