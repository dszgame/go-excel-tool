---@class ItemConfig
---@field nItemID int 道具id
---@field sName string 道具名称
---@field nType int 道具类型
---@field nUseType int 是否自动使用
---@field sContent string 道具描述
---@field tAttr table 附加属性
---@field tTestAuto auto 自动字段
---@field tAttr2 table 多维数组
---@return table<string, ItemConfig>
return {
	[1001]={nItemID=1001,sName=[[金疮药]],nType=1,sContent=[[这是一瓶金疮药，一口奶100HP]],tAttr={1,100,},tTestAuto=[[字符串测试]],tAttr2={{1,100,},{2,200,},},},
	[1002]={nItemID=1002,sName=[[大号金疮药]],nType=1,nUseType=1,sContent=[[这是一瓶大号金疮药，一口奶500HP]],tAttr={1,500,},tTestAuto=1,tAttr2={{1,100,},{2,200,},},},
}