package model

import (
	"encoding/json"
	"fmt"
	"time"
	"wallet_chain.com/global"
	"wallet_chain.com/utils"
)

type WalletChain struct {
	Id                int              `gorm:"column:id;primary_key;AUTO_INCREMENT"  form:"id"`
	Title             string           `gorm:"column:title" form:"title"`
	Icon              string           `gorm:"column:icon" json:"-"`
	Pid               int              `gorm:"column:pid" json:"-"`
	Chain             string           `gorm:"column:chain" form:"chain"`
	Label             string           `gorm:"column:label" form:"label"`
	CreateTime        time.Time        `gorm:"column:create_time" json:"-"`
	UpdateTime        time.Time        `gorm:"column:update_time" json:"-"`
	Remarks           string           `gorm:"column:remarks" json:"-"`
	Statuswithdraw    *int             `gorm:"column:statuswithdraw" form:"statuswithdraw"`
	Statusrecharge    *int             `gorm:"column:statusrecharge" form:"statusrecharge"`
	Address           *string          `gorm:"column:address;" form:"address"` // 地址
	IconUrl           string           `gorm:"-" `                             // 地址
	SubCategory       []*WalletChain   `gorm:"foreignKey:pid" json:"son"`
	WithdrawFee       *float64         `gorm:"withdraw_fee" form:"withdraw_fee"`
	InnerOrder        int              `gorm:"inner_order" json:"-"`
	MinWithdrawAmount *float64         `gorm:"min_withdraw_amount" form:"min_withdraw_amount"`
	MinRechargeAmount *float64         `gorm:"min_recharge_amount" form:"min_recharge_amount"`
	Type              int              `gorm:"type" form:"type"`
	Children          *WalletChainTree `json:"children" gorm:"-"`
}

func (m *WalletChain) TableName() string {
	return "wallet_chain"
}

type WalletChainTree []*WalletChain

func (m *WalletChain) MarshalJSON() ([]byte, error) {
	type Alias WalletChain // 创建一个新的类型别名以避免无限递归调用 MarshalJSON() 方法
	aux := struct {
		*Alias         // 将原始User的所有字段嵌入到aux中，但不包括Birthday（除非你已经处理了它）
		IconUrl string `json:"IconUrl"` // 自定义格式化时间字段名和格式化逻辑
		Name    string `json:"name"` // 自定义格式化时间字段名和格式化逻辑
	}{
		Alias:   (*Alias)(m),                                                      // 将原始User的所有公共字段赋值给aux的Alias部分（不包括Birthday）
		IconUrl: fmt.Sprintf("%s%s", global.SHOP_CONFIG.System.WebApiURL, m.Icon), // 使用自定义的时间格式化逻辑（例如：年-月-日）

	}
	return json.Marshal(aux) // 序列化aux结构体，其中包括格式化后的Birthday字段和所有其他原始字段。
}

func (*WalletChain) ToTree(data WalletChainTree, language string) WalletChainTree {
	// 定义 HashMap 的变量，并初始化
	TreeData := make(map[int]*WalletChain)
	// 先重组数据：以数据的ID作为外层的key编号，以便下面进行子树的数据组合
	for _, item := range data {
		TreeData[item.Id] = item
	}
	// 定义 RoleTrees 结构体
	var TreeDataList WalletChainTree
	// 开始生成树形
	for _, item := range TreeData {
		// 如果没有根节点树，则为根节点
		if item.Pid == 0 {
			// 追加到 TreeDataList 结构体中
			item.Title = utils.Languagebycode(language, item.Title)
			TreeDataList = append(TreeDataList, item)
			// 跳过该次循环
			continue
		}
		item.Title = utils.Languagebycode(language, item.Title)
		// 通过 上面的 TreeData HashMap的组合，进行判断是否存在根节点
		// 如果存在根节点，则对应该节点进行处理
		if pItem, ok := TreeData[item.Pid]; ok {
			fmt.Println(pItem.Id)
			// 判断当次循环是否存在子节点，如果没有则作为子节点进行组合
			if pItem.Children == nil {
				// 写入子节点
				children := WalletChainTree{item}
				// 插入到 当次结构体的子节点字段中，以指针的方式
				pItem.Children = &children
				pItem.Title = utils.Languagebycode(language, pItem.Title)
				// 跳过当前循环
				continue
			}

			// 以指针地址的形式进行追加到结构体中
			*pItem.Children = append(*pItem.Children, item)
		}

	}

	return TreeDataList
}
