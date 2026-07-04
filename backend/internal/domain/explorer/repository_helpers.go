package explorer

import "github.com/jackc/pgx/v5"

func scanNodes(rows pgx.Rows) ([]Node, error) {
	result := make([]Node, 0)
	for rows.Next() {
		var n Node
		if err := rows.Scan(
			&n.ID, &n.ParentKind, &n.ParentID, &n.NodeType, &n.DraftPointID, &n.Title,
			&n.OrderIndex, &n.Status, &n.CreatedAt, &n.UpdatedAt,
			&n.SourceType, &n.SourceURL, &n.ContentID, &n.Summary,
			&n.ResearchType, &n.LayoutType, &n.BlockCount,
		); err != nil {
			return nil, err
		}
		result = append(result, n)
	}
	return result, rows.Err()
}

func scanNodesWithLearningStatus(rows pgx.Rows) ([]Node, error) {
	result := make([]Node, 0)
	for rows.Next() {
		var n Node
		if err := rows.Scan(
			&n.ID, &n.ParentKind, &n.ParentID, &n.NodeType, &n.DraftPointID, &n.LearningStatus,
			&n.Title, &n.OrderIndex, &n.Status, &n.CreatedAt, &n.UpdatedAt,
			&n.SourceType, &n.SourceURL, &n.ContentID, &n.Summary,
			&n.ResearchType, &n.LayoutType, &n.BlockCount,
		); err != nil {
			return nil, err
		}
		result = append(result, n)
	}
	return result, rows.Err()
}
